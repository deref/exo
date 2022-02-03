package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/deref/exo/internal/util/mathutil"
)

type EventResolver struct {
	Q *RootResolver
	EventRow
}

type EventRow struct {
	ULID        ULID    `db:"ulid"`
	Type        string  `db:"type"`
	Message     string  `db:"message"`
	Tags        Tags    `db:"tags"`
	WorkspaceID *string `db:"workspace_id"`
	StackID     *string `db:"stack_id"`
	ComponentID *string `db:"component_id"`
	JobID       *string `db:"job_id"`
	TaskID      *string `db:"task_id"`
}

func (r *EventRow) ID() string {
	return r.ULID.String()
}

type createEventArgs struct {
	Source    StreamSourceResolver
	Timestamp *Instant
	Type      string
	Message   string
	Tags      *Tags
}

func (r *MutationResolver) CreateEvent(ctx context.Context, args createEventArgs) (*EventResolver, error) {
	row := EventRow{
		ULID:    r.mustNextULID(ctx),
		Type:    args.Type,
		Message: args.Message,
		XXX:     SET_RELATED_IDS,
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

func (r *EventResolver) Timestamp() Instant {
	return r.ULID.Timestamp()
}

type eventQuery struct {
	Filter eventFilter
	Cursor string
	Prev   int
	Next   int
}

type eventFilter struct {
	WorkspaceID string
	StackID     string
	ComponentID string
	JobID       string
	TaskID      string
	IContains   string
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
		cursor, err = r.latestEventCursor(ctx, q.Filter)
		if err != nil {
			return nil, fmt.Errorf("getting latest cursor: %w", err)
		}
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
			WHERE (
					 (? IS NOT NULL AND workspace_id = ?)
				OR (? IS NOT NULL AND stack_id = ?)
				OR (? IS NOT NULL AND component_id = ?)
				OR (? IS NOT NULL AND job_id = ?)
				OR (? IS NOT NULL AND task_id = ?)
			)
			AND instr(lower(message), ?) <> 0
			AND ulid < ?
			ORDER BY ulid DESC
			LIMIT ?
		`
	} else {
		query = `
			SELECT *
			FROM event
			WHERE (
					 (? IS NOT NULL AND workspace_id = ?)
				OR (? IS NOT NULL AND stack_id = ?)
				OR (? IS NOT NULL AND component_id = ?)
				OR (? IS NOT NULL AND job_id = ?)
				OR (? IS NOT NULL AND task_id = ?)
			)
			AND instr(lower(message), ?) <> 0
			AND ? < ulid
			ORDER BY ulid ASC
			LIMIT ?
		`
	}
	var rows []EventRow
	filter := q.Filter
	if err := r.DB.SelectContext(ctx, &rows, query,
		filter.WorkspaceID, filter.WorkspaceID,
		filter.StackID, filter.StackID,
		filter.ComponentID, filter.ComponentID,
		filter.JobID, filter.JobID,
		filter.TaskID, filter.TaskID,
		filter.IContains,
		cursor, limit,
	); err != nil {
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

func (r *QueryResolver) latestEventCursor(ctx context.Context, filter eventFilter) (string, error) {
	event, err := r.latestEvent(ctx, filter)
	if event == nil || err != nil {
		return "", err
	}
	return incrementEventCursor(event.ID()), nil
}

func (r *QueryResolver) latestEvent(ctx context.Context, filter eventFilter) (*EventResolver, error) {
	var row EventRow
	err := r.DB.GetContext(ctx, &row, `
		SELECT *
		FROM event
		WHERE (
			   (? IS NOT NULL AND workspace_id = ?)
			OR (? IS NOT NULL AND stack_id = ?)
			OR (? IS NOT NULL AND component_id = ?)
			OR (? IS NOT NULL AND job_id = ?)
			OR (? IS NOT NULL AND task_id = ?)
		)
		AND instr(lower(message), ?) <> 0
		ORDER BY ulid DESC
		LIMIT 1
	`,
		filter.WorkspaceID, filter.WorkspaceID,
		filter.StackID, filter.StackID,
		filter.ComponentID, filter.ComponentID,
		filter.JobID, filter.JobID,
		filter.TaskID, filter.TaskID,
		filter.IContains,
	)
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

type eventSubscription struct {
	Filter eventFilter
	Cursor string
}

func (r *SubscriptionResolver) events(ctx context.Context, sub eventSubscription) (<-chan *EventResolver, error) {
	logger := r.SystemLog.Sublogger("events subscription")
	c := make(chan *EventResolver)
	go func() {
		defer close(c)

		// Poll for events.
		cursor := sub.Cursor
		for {
			page, err := r.findEvents(ctx, eventQuery{
				RelevantTo: sub.RelevantTo,
				Cursor:     cursor,
			})
			if err != nil {
				logger.Infof("error finding events: %v", err)
				return
			}
			events := page.Items
			cursor = page.NextCursor

			// Emit events.
			for _, event := range events {
				select {
				case <-ctx.Done():
					return
				case c <- event:
				}
			}

			// Allow time for additional events to occur.
			if len(events) == 0 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(30 * time.Millisecond):
				}
			}
		}
	}()
	return c, nil
}

func (r *EventResolver) Workspace(ctx context.Context) (*WorkspaceResolver, error) {
	return r.Q.workspaceByID(ctx, r.WorkspaceID)
}

func (r *EventResolver) Stack(ctx context.Context) (*StackResolver, error) {
	return r.Q.stackByID(ctx, r.StackID)
}

func (r *EventResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByID(ctx, r.ComponentID)
}

func (r *EventResolver) Job() *JobResolver {
	return r.Q.jobByID(r.JobID)
}

func (r *EventResolver) Task(ctx context.Context) (*TaskResolver, error) {
	return r.Q.taskByID(ctx, r.TaskID)
}
