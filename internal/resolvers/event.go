package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/util/mathutil"
	"github.com/jmoiron/sqlx"
)

type EventResolver struct {
	Q *RootResolver
	EventRow
}

type EventRow struct {
	ULID     ULID   `db:"ulid"`
	StreamID string `db:"stream_id"`
	Message  string `db:"message"`
	Tags     Tags   `db:"tags"`
}

func (r *EventResolver) ID() string {
	return r.ULID.String()
}

func (r *MutationResolver) CreateEvent(ctx context.Context, args struct {
	StreamID  string
	Timestamp *Instant
	Message   string
	Tags      *Tags
}) (*EventResolver, error) {
	row := EventRow{
		StreamID: args.StreamID,
		ULID:     r.mustNextULID(ctx),
		Message:  args.Message,
	}
	if args.Tags != nil {
		row.Tags = *args.Tags
	}
	if err := r.insertRow(ctx, "event", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &EventResolver{
		Q:        r,
		EventRow: row,
	}, nil
}

func (r *EventResolver) Stream(ctx context.Context) (*StreamResolver, error) {
	return r.Q.streamById(ctx, &r.StreamID)
}

func (r *EventResolver) Timestamp() Instant {
	return r.ULID.Timestamp()
}

type eventQuery struct {
	StreamIDs []string
	Cursor    string
	Prev      int
	Next      int
	Filter    string
}

type EventPageResolver struct {
	Items      []*EventResolver
	PrevCursor string
	NextCursor string
}

const defaultEventPageSize = 500
const maxEventPageSize = 10000

func (r *QueryResolver) findEvents(ctx context.Context, q eventQuery) (*EventPageResolver, error) {
	output := &EventPageResolver{}

	cursor := q.Cursor
	if q.Cursor == "" {
		var err error
		cursor, err = r.latestEventCursor(ctx, q.StreamIDs)
		if err != nil {
			return nil, fmt.Errorf("getting latest cursor: %w", err)
		}
	}

	if len(q.StreamIDs) == 0 {
		output.PrevCursor = cursor
		output.NextCursor = cursor
		return output, nil
	}

	limit := defaultEventPageSize
	reverse := false
	if q.Next > 0 {
		if q.Prev > 0 {
			return nil, errors.New("only one of prev or next may be specified")
		}
		limit = q.Next
	} else if q.Prev > 0 {
		limit = q.Prev
		reverse = true
	} else {
		// Default limit and forward.
	}
	limit = mathutil.IntClamp(limit, 0, maxEventPageSize)

	var query string
	if reverse {
		query = `
			SELECT *
			FROM event
			WHERE stream_id IN (?)
			AND id < ?
			AND instr(lower(message), ?) <> 0
			ORDER BY id DESC
			LIMIT ?
		`
	} else {
		query = `
			SELECT *
			FROM event
			WHERE stream_id IN (?)
			AND ? < id
			AND instr(lower(message), ?) <> 0
			ORDER BY id ASC
			LIMIT ?
		`
	}
	query, args := mustSqlIn(query, q.StreamIDs, cursor, q.Filter, limit)

	var rows []EventRow
	if err := r.DB.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("querying: %w", err)
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
		output.PrevCursor = rows[0].ULID.String()
		output.NextCursor = incrementEventCursor(rows[len(output.Items)-1].ULID.String())
	}

	output.Items = make([]*EventResolver, len(rows))
	for i, row := range rows {
		output.Items[i] = &EventResolver{
			Q:        r,
			EventRow: row,
		}
	}
	return output, nil
}

func (r *QueryResolver) latestEventCursor(ctx context.Context, streamIDs []string) (string, error) {
	if len(streamIDs) == 0 {
		return r.mustNextULID(ctx).String(), nil
	}
	query, args, err := sqlx.In(`
		SELECT COALESCE(MAX(id), "")
		FROM event
		WHERE stream IN (?)
	`, streamIDs)
	if err != nil {
		panic(err)
	}
	var id string
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	default:
		return incrementEventCursor(id), nil
	}
}

func incrementEventCursor(id string) string {
	return id + "0"
}
