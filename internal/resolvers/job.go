package resolvers

import (
	"context"
	"fmt"

	. "github.com/deref/exo/internal/scalars"
)

type JobResolver struct {
	Q *RootResolver
	// Jobs are synthetic entities with an ID equal to that of their root task.
	ID string
}

func (r *QueryResolver) jobByID(id *string) *JobResolver {
	if id == nil {
		return nil
	}
	return &JobResolver{
		Q:  r,
		ID: *id,
	}
}

func (r *QueryResolver) jobByRootTaskID(rootTaskID *string) *JobResolver {
	return r.jobByID(rootTaskID)
}

func (r *JobResolver) URL() string {
	return r.Q.Routes.jobURL(r.ID)
}

func (r *JobResolver) RootTask(ctx context.Context) (*TaskResolver, error) {
	task, err := r.Q.taskByID(ctx, &r.ID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("cannot find root task: %w", err)
	}
	return task, nil
}

func (r *JobResolver) Tasks(ctx context.Context) ([]*TaskResolver, error) {
	return r.Q.tasksByJobID(ctx, r.ID)
}

func (r *MutationResolver) CancelJob(ctx context.Context, args struct {
	ID string
}) error {
	return r.cancelJob(ctx, args.ID)
}

func (r *MutationResolver) cancelJob(ctx context.Context, id string) error {
	now := Now(ctx)
	_, err := r.DB.ExecContext(ctx, `
		UPDATE task
		SET canceled = COALESCE(canceled, ?)
		WHERE job_id = ?
	`, now, id)
	return err
}

func (r *SubscriptionResolver) WatchJob(ctx context.Context, args struct {
	ID    string
	After *ULID
}) (<-chan *EventResolver, error) {
	jobID := args.ID
	job := r.jobByID(&jobID)

	// Subscribe to events.
	filter := eventFilter{
		JobID: jobID,
	}
	if args.After != nil {
		filter.After = *args.After
	}
	events, err := r.events(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Create output channel with synthetic initial event.
	c := make(chan *EventResolver, 1)
	{
		watched, err := job.eventPrototype(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving initial event: %w", err)
		}
		watched.ULID = r.mustNextULID(ctx)
		watched.Type = "JobWatched"
		c <- &EventResolver{
			Q:        r,
			EventRow: watched,
		}
	}

	// Pipe events to output.
	go func() {
		defer close(c)
		for event := range events {
			select {
			case c <- event:
			case <-ctx.Done():
				return
			}
		}
	}()

	return c, nil
}

func (r *JobResolver) eventPrototype(ctx context.Context) (row EventRow, err error) {
	task, err := r.RootTask(ctx)
	if err != nil {
		return EventRow{}, fmt.Errorf("resolving root task: %w", err)
	}
	return task.eventPrototype(ctx)
}
