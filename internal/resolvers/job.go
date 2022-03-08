package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/scalars"
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

func (r *MutationResolver) CreateJob(ctx context.Context, args struct {
	Mutation  string
	Arguments JSONObject
}) (*JobResolver, error) {
	return r.createJob(ctx, args.Mutation, args.Arguments)
}

func (r *MutationResolver) createJob(ctx context.Context, mutation string, arguments map[string]interface{}) (*JobResolver, error) {
	task, err := r.createRootTask(ctx, mutation, arguments)
	if err != nil {
		return nil, err
	}
	return &JobResolver{
		Q:  r,
		ID: task.ID,
	}, nil
}

func (r *JobResolver) URL() string {
	return r.Q.Routes().jobURL(r.ID)
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

// See also cancelTask and cancelSubtasks.
func (r *MutationResolver) cancelJob(ctx context.Context, id string) error {
	now := Now(ctx)
	_, err := r.db.ExecContext(ctx, `
		UPDATE task
		SET canceled = COALESCE(canceled, ?)
		WHERE job_id = ?
		AND finished IS NULL
	`, now, id)
	return err
}

func (r *SubscriptionResolver) WatchJob(ctx context.Context, args struct {
	ID    string
	After *ULID
	Debug *bool
}) (<-chan *EventResolver, error) {
	jobID := args.ID
	job := r.jobByID(&jobID)
	rootTask, err := job.RootTask(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving root task: %w", err)
	}

	// Subscribe to events.
	filter := eventFilter{
		JobID: jobID,
	}
	if args.After == nil {
		filter.After = scalars.InstantToULID(rootTask.Created)
	} else {
		filter.After = *args.After
	}
	if args.Debug != nil && *args.Debug {
		filter.System = true
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

func (r *JobResolver) Stream() *StreamResolver {
	return r.Q.streamForSource("Job", r.ID)
}

func (r *JobResolver) eventPrototype(ctx context.Context) (row EventRow, err error) {
	task, err := r.RootTask(ctx)
	if err != nil {
		return EventRow{}, fmt.Errorf("resolving root task: %w", err)
	}
	prototype, err := task.eventPrototype(ctx)
	prototype.SourceType = "Job"
	return prototype, err
}

func (r *JobResolver) Updated(ctx context.Context) (Instant, error) {
	var res Instant
	if err := r.Q.db.GetContext(ctx, &res, `
		SELECT MAX(updated)
		FROM task
		WHERE job_id = ?
	`, r.ID,
	); err != nil {
		return Instant{}, err
	}
	return res, nil
}

func (r *QueryResolver) isJobCompleted(ctx context.Context, jobID string) (bool, error) {
	return r.isTaskCompleted(ctx, jobID)
}
