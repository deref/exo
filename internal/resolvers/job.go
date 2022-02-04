package resolvers

import (
	"context"

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
	return r.Q.taskByID(ctx, &r.ID)
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
}) (<-chan *WatchJobOutput, error) {
	jobID := args.ID
	job := r.jobByID(&jobID)

	// Subscribe to events.
	filter := eventFilter{
		JobID: jobID,
	}
	if args.After != nil {
		filter.After = *args.After
	}
	eventC, err := r.events(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Create output channel with initial output.
	outputC := make(chan *WatchJobOutput, 1)
	outputC <- &WatchJobOutput{
		Job: job,
	}

	// Pipe events to output.
	go func() {
		defer close(outputC)
		for event := range eventC {
			select {
			case outputC <- &WatchJobOutput{
				Job:   job,
				Event: event,
			}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return outputC, nil
}

type WatchJobOutput struct {
	Job   *JobResolver
	Event *EventResolver
}
