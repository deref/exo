package resolvers

import "context"

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

func (r *JobResolver) Task(ctx context.Context) (*TaskResolver, error) {
	return r.Q.taskByID(ctx, &r.ID)
}
