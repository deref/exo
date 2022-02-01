package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/util/mathutil"
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

func (r *EventRow) ID() string {
	return r.ULID.String()
}

type createEventArgs struct {
	StreamID  string
	Timestamp *Instant
	Message   string
	Tags      *Tags
}

func (r *MutationResolver) CreateEvent(ctx context.Context, args createEventArgs) (*EventResolver, error) {
	row := EventRow{
		StreamID: args.StreamID,
		ULID:     r.mustNextULID(ctx),
		Message:  args.Message,
	}
	if args.Tags == nil {
		row.Tags = make(Tags)
	} else {
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
			AND ulid < ?
			AND instr(lower(message), ?) <> 0
			ORDER BY ulid DESC
			LIMIT ?
		`
	} else {
		query = `
			SELECT *
			FROM event
			WHERE stream_id IN (?)
			AND ? < ulid
			AND instr(lower(message), ?) <> 0
			ORDER BY ulid ASC
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
		output.NextCursor = incrementEventCursor(rows[len(output.Items)-1].ID())
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
	event, err := r.latestEvent(ctx, streamIDs)
	if event == nil || err != nil {
		return "", err
	}
	return incrementEventCursor(event.ID()), nil
}

func (r *QueryResolver) latestEvent(ctx context.Context, streamIDs []string) (*EventResolver, error) {
	q, args := mustSqlIn(`
		SELECT *
		FROM event
		WHERE stream_id IN (?)
		ORDER BY ulid DESC
		LIMIT 1
	`, streamIDs)
	var row EventRow
	err := r.DB.GetContext(ctx, &row, q, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &EventResolver{
		Q:        r,
		EventRow: row,
	}, nil
}

func incrementEventCursor(id string) string {
	return id + "0"
}
