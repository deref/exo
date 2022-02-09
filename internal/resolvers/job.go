package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/deref/exo/internal/chrono"
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

	prototype, err := job.eventPrototype(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving event prototype: %w", err)
	}

	// Create output channel with synthetic initial event.
	c := make(chan *EventResolver, 1)
	{
		jobWatchedEvent := prototype
		jobWatchedEvent.Type = "JobWatched"
		c <- r.newSyntheticEvent(ctx, jobWatchedEvent)
	}

	t := chrono.Now(ctx)

	// Pipe events to output. Also periodically emits synthetic events for
	// reporting progress without having to have a physical event for every task
	// progress update.
	go func() {
		defer close(c)

		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()

		for {
			var event *EventResolver
			select {
			case <-ctx.Done():
				return
			case event = <-events:
			case <-ticker.C:
				// Check for progress.
				updated, err := job.Updated(ctx)
				if err != nil {
					r.SystemLog.Infof("resolving job updated time: %v", err)
					break
				}
				if !updated.GoTime().After(t) {
					continue
				}
				t = updated.GoTime()

				// Emit periodic progress event.
				jobUpdatedEvent := prototype
				jobUpdatedEvent.Type = "JobUpdated"
				event = r.newSyntheticEvent(ctx, jobUpdatedEvent)
			}

			// Forward event.
			select {
			case <-ctx.Done():
				return
			case c <- event:
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

func (r *JobResolver) Updated(ctx context.Context) (Instant, error) {
	var res Instant
	if err := r.Q.DB.GetContext(ctx, &res, `
		SELECT MAX(updated)
		FROM task
		WHERE job_id = ?
	`, r.ID,
	); err != nil {
		return Instant{}, err
	}
	return res, nil
}