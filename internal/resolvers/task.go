package resolvers

import "context"

type TaskResolver struct {
	Q *RootResolver
	TaskRow
}

type TaskRow struct {
	ID              string   `db:"id"`
	JobID           string   `db:"job_id"`
	ParentID        *string  `db:"parent_id"`
	Mutation        string   `db:"mutation"`
	Variables       string   `db:"variables"`
	WorkerID        *string  `db:"worker_id"`
	Status          string   `db:"status"`
	Created         Instant  `db:"created"`
	Updated         Instant  `db:"updated"`
	Started         *Instant `db:"started"`
	Finished        *Instant `db:"finished"`
	ProgressCurrent *int32   `db:"progress_current"`
	ProgressTotal   *int32   `db:"progress_total"`
	Message         *string  `db:"message"`
}

func (r *QueryResolver) taskByID(ctx context.Context, id *string) (*TaskResolver, error) {
	t := &TaskResolver{}
	err := r.getRowByID(ctx, &t.TaskRow, `
		SELECT
			id,
			job_id,
			parent_id,
			mutation,
			variables,
			worker_id,
			status,
			created,
			updated,
			started,
			finished,
			progress_current,
			progress_total,
			message
		FROM task
		WHERE id = ?
	`, id)
	if t.ID == "" {
		t = nil
	}
	return t, err
}

func (r *QueryResolver) TasksByJobID(ctx context.Context, args struct {
	JobID string
}) ([]*TaskResolver, error) {
	return r.tasksByJobID(ctx, args.JobID)
}

func (r *QueryResolver) tasksByJobID(ctx context.Context, jobID string) ([]*TaskResolver, error) {
	var rows []TaskRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT
			id,
			job_id,
			parent_id,
			mutation,
			variables,
			worker_id,
			status,
			created,
			updated,
			started,
			finished,
			progress_current,
			progress_total,
			message
		FROM task
		WHERE job_id = ?
		ORDER BY task.id ASC
  `, jobID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*TaskResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &TaskResolver{
			Q:       r,
			TaskRow: row,
		}
	}
	return resolvers, nil
}

func (r *TaskResolver) Job(ctx context.Context) (*TaskResolver, error) {
	return r.Q.taskByID(ctx, &r.JobID)
}

func (r *TaskResolver) Parent(ctx context.Context) (*TaskResolver, error) {
	return r.Q.taskByID(ctx, r.ParentID)
}

func (r *TaskResolver) Label() string {
	switch r.Mutation {
	default:
		return r.Mutation
	}
}

func (r *TaskResolver) Progress() (*ProgressResolver, error) {
	if r.ProgressCurrent == nil || r.ProgressTotal == nil {
		return nil, nil
	}
	return &ProgressResolver{
		Current: *r.ProgressCurrent,
		Total:   *r.ProgressTotal,
	}, nil
}
