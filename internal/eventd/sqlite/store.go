package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/eventd/api"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/mathutil"
)

type Store struct {
	DB    *sqlx.DB
	IDGen *gensym.ULIDGenerator
}

func (sto *Store) ClearEvents(ctx context.Context, input *api.ClearEventsInput) (*api.ClearEventsOutput, error) {
	if len(input.Streams) > 0 {
		query, args, err := sqlx.In(`
		DELETE FROM event
		WHERE stream IN (?)
	`, input.Streams)
		if err != nil {
			panic(err)
		}
		if _, err := sto.DB.ExecContext(ctx, query, args...); err != nil {
			return nil, err
		}
	}
	return &api.ClearEventsOutput{}, nil
}

func (sto *Store) DescribeStreams(ctx context.Context, input *api.DescribeStreamsInput) (*api.DescribeStreamsOutput, error) {
	output := api.DescribeStreamsOutput{
		Streams: []api.StreamDescription{},
	}
	if len(input.Names) > 0 {
		query, args, err := sqlx.In(`
			SELECT stream, max(timestamp) AS last_event_at
			FROM event
			WHERE name IN (?)
			ORDER BY stream ASC
			GROUP BY stream
		`, input.Names)
		if err != nil {
			panic(err)
		}
		rows, err := sto.DB.QueryxContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("querying: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var stream api.StreamDescription
			var lastEventAtNano int64
			if err := rows.Scan(&stream.Name, &lastEventAtNano); err != nil {
				return nil, fmt.Errorf("scanning: %w", err)
			}
			lastEventAtIso := chrono.NanoToIso(lastEventAtNano)
			stream.LastEventAt = &lastEventAtIso
			output.Streams = append(output.Streams, stream)
		}
		if rows.Err() != nil {
			return nil, fmt.Errorf("advancing rows: %w", rows.Err())
		}
	}
	return &output, nil
}

func (sto *Store) AddEvent(ctx context.Context, input *api.AddEventInput) (*api.AddEventOutput, error) {
	if input.Stream == "" {
		return nil, errutil.NewHTTPError(http.StatusBadRequest, "stream is required")
	}

	id := sto.nextID(ctx)

	timestamp, err := chrono.ParseIsoToNano(input.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("parsing timestamp: %w", err)
	}

	if _, err := sto.DB.ExecContext(ctx, `
		INSERT INTO event ( stream, id, timestamp, message )
		VALUES ( ?, ?, ?, ? )
	`, input.Stream, id, timestamp, input.Message); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &api.AddEventOutput{}, nil
}

// Generate an id that is guaranteed to be monotonically increasing within this process.
func (sto *Store) nextID(ctx context.Context) string {
	lid, err := sto.IDGen.NextID(ctx)
	if err != nil {
		panic(err)
	}
	lidBytes, err := lid.MarshalText()
	if err != nil {
		panic(err)
	}
	return string(lidBytes)
}

const (
	defaultNextLimit = 500
	maxLimit         = 10000
)

func (sto *Store) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	output := api.GetEventsOutput{
		Items: []api.Event{},
	}

	var cursor string
	if input.Cursor == nil {
		var err error
		cursor, err = sto.getLatestCursor(ctx, input.Streams)
		if err != nil {
			return nil, fmt.Errorf("getting latest cursor: %w", err)
		}
	} else {
		cursor = *input.Cursor
	}

	if len(input.Streams) == 0 {
		output.PrevCursor = cursor
		output.NextCursor = cursor
		return &output, nil
	}

	limit := defaultNextLimit
	reverse := false
	if input.Next != nil {
		if input.Prev != nil {
			return nil, errors.New("only one of prev or next may be specified")
		}
		limit = *input.Next
	} else if input.Prev != nil {
		limit = *input.Prev
		reverse = true
	} else {
		// Default limit and forward.
	}
	limit = mathutil.IntClamp(limit, 0, maxLimit)

	var query string
	if reverse {
		query = `
			SELECT stream, id, timestamp, message
			FROM event
			WHERE stream IN (?)
			AND id < ?
			AND instr(lower(message), ?) <> 0
			ORDER BY id DESC
			LIMIT ?
		`
	} else {
		query = `
			SELECT stream, id, timestamp, message
			FROM event
			WHERE stream IN (?)
			AND ? < id
			AND instr(lower(message), ?) <> 0
			ORDER BY id ASC
			LIMIT ?
		`
	}
	query, args, err := sqlx.In(query, input.Streams, cursor, input.FilterStr, limit)
	if err != nil {
		panic(err)
	}

	rows, err := sto.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event api.Event
		var timestampNano int64
		if err := rows.Scan(&event.Stream, &event.ID, &timestampNano, &event.Message); err != nil {
			return nil, fmt.Errorf("scanning: %w", err)
		}
		event.Timestamp = chrono.NanoToIso(timestampNano)
		output.Items = append(output.Items, event)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("advancing rows: %w", rows.Err())
	}

	if reverse {
		l := 0
		r := len(output.Items) - 1
		for l < r {
			tmp := output.Items[l]
			output.Items[l] = output.Items[r]
			output.Items[r] = tmp
			l++
			r--
		}
	}

	output.PrevCursor = cursor
	output.NextCursor = cursor
	if len(output.Items) > 0 {
		output.PrevCursor = output.Items[0].ID
		output.NextCursor = incrementCursor(output.Items[len(output.Items)-1].ID)
	}

	return &output, nil
}

func (sto *Store) getLatestCursor(ctx context.Context, streams []string) (string, error) {
	if len(streams) == 0 {
		return sto.nextID(ctx), nil
	}
	query, args, err := sqlx.In(`
		SELECT COALESCE(MAX(id), "")
		FROM event
		WHERE stream IN (?)
	`, streams)
	if err != nil {
		panic(err)
	}
	var id string
	err = sto.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	default:
		return incrementCursor(id), nil
	}
}

func incrementCursor(id string) string {
	return id + "0"
}

func (sto *Store) RemoveOldEvents(ctx context.Context, input *api.RemoveOldEventsInput) (*api.RemoveOldEventsOutput, error) {
	// This is an inefficent way to keep only the most recent rows.
	// TODO: Use SQLITE_ENABLE_UPDATE_DELETE_LIMIT when go-sqlite supports it.
	// See <https://github.com/mattn/go-sqlite3/issues/787>.
	const maxEvents = 10000
	_, err := sto.DB.ExecContext(ctx, `
		DELETE FROM event
		WHERE id NOT IN (
			SELECT id
			FROM event
			ORDER BY id DESC
			LIMIT ?
		)
	`, maxEvents)
	if err != nil {
		return nil, err
	}
	return &api.RemoveOldEventsOutput{}, nil
}
