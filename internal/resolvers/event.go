package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/deref/exo/internal/api"
	. "github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/mathutil"
)

type EventResolver struct {
	Q *RootResolver
	EventRow
}

type EventRow struct {
	ULID        ULID       `db:"ulid"`
	Type        string     `db:"type"`
	Message     string     `db:"message"`
	Tags        JSONObject `db:"tags"`
	SourceType  string     `db:"source_type"`
	WorkspaceID *string    `db:"workspace_id"`
	StackID     *string    `db:"stack_id"`
	ComponentID *string    `db:"component_id"`
	JobID       *string    `db:"job_id"`
	TaskID      *string    `db:"task_id"`
}

func (r *EventRow) ID() string {
	return r.ULID.String()
}

func (r *MutationResolver) CreateEvent(ctx context.Context, args struct {
	SourceType string
	SourceID   string
	Type       string
	Message    string
}) (*EventResolver, error) {
	source, err := r.findEventSource(ctx, args.SourceType, args.SourceID)
	if err != nil {
		return nil, fmt.Errorf("resolving source: %w", err)
	}
	if source == nil {
		return nil, fmt.Errorf("cannot find event source: type=%q id=%q", args.SourceType, args.SourceID)
	}
	return r.createEvent(ctx, source, args.Type, args.Message)
}

func (r *MutationResolver) createEvent(ctx context.Context, source StreamSourceResolver, typ string, message string) (*EventResolver, error) {
	prototype, err := source.eventPrototype(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving event prototype: %w", err)
	}
	prototype.Type = typ
	prototype.Message = message
	return r.createEventFromPrototype(ctx, prototype)
}

func (r *MutationResolver) createEventFromPrototype(ctx context.Context, prototype EventRow) (*EventResolver, error) {
	row := prototype
	row.ULID = r.mustNextULID(ctx)
	if row.Type == "" {
		return nil, fmt.Errorf("invalid event type: %q", row.Type)
	}
	if row.SourceType == "" {
		return nil, errors.New("source type is required")
	}
	if row.Tags == nil {
		row.Tags = make(JSONObject)
	}
	if err := r.insertRow(ctx, "event", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &EventResolver{
		Q:        r,
		EventRow: row,
	}, nil
}

func (r *MutationResolver) newSyntheticEvent(ctx context.Context, row EventRow) *EventResolver {
	row.ULID = r.mustNextULID(ctx)
	return &EventResolver{
		Q:        r,
		EventRow: row,
	}
}

func (r *MutationResolver) mustNextULID(ctx context.Context) api.ULID {
	res, err := r.ulidgen.NextID(ctx)
	if err != nil {
		panic(err)
	}
	return ULID(res)
}

func (r *EventRow) Timestamp() Instant {
	return r.ULID.Timestamp()
}

type eventQuery struct {
	Filter eventFilter
	Cursor ULID
	Prev   int
	Next   int
}

type eventFilter struct {
	Before      ULID
	After       ULID
	System      bool
	WorkspaceID string
	StackID     string
	// TODO: If ComponentID and StackID are both set, probably want to remove
	// events from encapsulated components by checking component.parent_id is null.
	ComponentID string
	JobID       string
	TaskID      string
	IContains   string
}

type EventPageResolver struct {
	Items      []*EventResolver
	PrevCursor ULID
	NextCursor ULID
}

const defaultEventPageSize = 500
const maxEventPageSize = 10000

func (r *QueryResolver) findEvents(ctx context.Context, q eventQuery) (*EventPageResolver, error) {
	cursor := q.Cursor
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

	filter := q.Filter
	after := filter.After
	before := filter.Before
	if before == (ULID{}) {
		before = InfiniteULID
	}

	var query string
	if reverse {
		query = `
			SELECT *
			FROM event
			WHERE (
				(? AND source_type == 'System')
				OR (? != '' AND workspace_id = ?)
				OR (? != '' AND stack_id = ?)
				OR (? != '' AND component_id = ?)
				OR (? != '' AND job_id = ?)
				OR (? != '' AND task_id = ?)
			)
			AND instr(lower(message), ?) <> 0
			AND ulid BETWEEN ? AND ?
			ORDER BY ulid DESC
			LIMIT ?
		`
		before = ULIDMax(before, cursor)
	} else {
		query = `
			SELECT *
			FROM event
			WHERE (
				(? AND source_type == 'System')
				OR (? != '' AND workspace_id = ?)
				OR (? != '' AND stack_id = ?)
				OR (? != '' AND component_id = ?)
				OR (? != '' AND job_id = ?)
				OR (? != '' AND task_id = ?)
			)
			AND instr(lower(message), ?) <> 0
			AND ulid BETWEEN ? AND ?
			ORDER BY ulid ASC
			LIMIT ?
		`
		after = ULIDMax(after, cursor)
	}
	var rows []EventRow
	if err := r.db.SelectContext(ctx, &rows, query,
		filter.System,
		filter.WorkspaceID, filter.WorkspaceID,
		filter.StackID, filter.StackID,
		filter.ComponentID, filter.ComponentID,
		filter.JobID, filter.JobID,
		filter.TaskID, filter.TaskID,
		filter.IContains,
		after.String(), before.String(), limit,
	); err != nil {
		return nil, fmt.Errorf("querying: %w", err)
	}

	if reverse {
		l := 0
		r := len(rows) - 1
		for l < r {
			tmp := rows[l]
			rows[l] = rows[r]
			rows[r] = tmp
			l++
			r--
		}
	}

	output := &EventPageResolver{
		Items:      eventRowsToResolvers(r, rows),
		PrevCursor: cursor,
		NextCursor: cursor,
	}

	if len(rows) > 0 {
		output.PrevCursor = rows[0].ULID
		output.NextCursor = IncrementULID(rows[len(rows)-1].ULID)
	}

	return output, nil
}

func eventRowsToResolvers(r *RootResolver, rows []EventRow) []*EventResolver {
	resolvers := make([]*EventResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &EventResolver{
			Q:        r,
			EventRow: row,
		}
	}
	return resolvers
}

// XXX When if ever is it appropriate to use this?
func (r *QueryResolver) latestEventCursor(ctx context.Context, filter eventFilter) (ULID, error) {
	event, err := r.latestEvent(ctx, filter)
	if event == nil || err != nil {
		return ULID{}, err
	}
	return IncrementULID(event.ULID), nil
}

func (r *QueryResolver) latestEvent(ctx context.Context, filter eventFilter) (*EventResolver, error) {
	var row EventRow
	err := r.db.GetContext(ctx, &row, `
		SELECT *
		FROM event
		WHERE (
			(? AND source_type == 'System')
			OR (? != '' AND workspace_id = ?)
			OR (? != '' AND stack_id = ?)
			OR (? != '' AND component_id = ?)
			OR (? != '' AND job_id = ?)
			OR (? != '' AND task_id = ?)
		)
		AND instr(lower(message), ?) <> 0
		ORDER BY ulid DESC
		LIMIT 1
	`,
		filter.System,
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
	Cursor ULID
}

func (r *SubscriptionResolver) events(ctx context.Context, filter eventFilter) (<-chan *EventResolver, error) {
	logger := r.SystemLog.Sublogger("events subscription")
	c := make(chan *EventResolver)
	go func() {
		defer close(c)

		// Poll for events.
		q := eventQuery{
			Filter: filter,
		}
		for {
			page, err := r.findEvents(ctx, q)
			if err != nil {
				logger.Infof("error finding events: %v", err)
				return
			}
			events := page.Items
			q.Cursor = page.NextCursor

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

func (r *EventResolver) SourceID() string {
	switch r.SourceType {
	case "System":
		return "SYSTEM"
	case "Job":
		return *r.JobID
	case "Task":
		return *r.TaskID
	default:
		panic(fmt.Errorf("unexpected source type: %q", r.SourceType))
	}
}

func (r *EventResolver) Stream() *StreamResolver {
	return r.Q.streamForSource(r.SourceType, r.SourceID())
}

func (r *EventResolver) Source(ctx context.Context) (source StreamSourceResolver, err error) {
	return r.Stream().Source(ctx)
}
