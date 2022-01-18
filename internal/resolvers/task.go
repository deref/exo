package resolvers

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/gensym"
)

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

const TaskStatusPending = "pending"
const TaskStatusRunning = "running"
const TaskStatusSuccess = "success"
const TaskStatusFailure = "failure"

func (r *MutationResolver) NewTask(ctx context.Context, args struct {
	ParentID  *string
	Mutation  string
	Variables string
}) (*TaskResolver, error) {
	id := gensym.RandomBase32()
	now := Now(ctx)
	row := TaskRow{
		ID:        id,
		ParentID:  args.ParentID,
		Mutation:  args.Mutation,
		Variables: args.Variables,
		Status:    TaskStatusPending,
		Created:   now,
		Updated:   now,
	}
	if args.ParentID == nil {
		row.JobID = id
	} else {
		parent, err := r.taskByID(ctx, args.ParentID)
		if err != nil {
			return nil, fmt.Errorf("resolving parent: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("no such parent: %q", *args.ParentID)
		}
		row.JobID = parent.JobID
	}
	if _, err := r.DB.ExecContext(ctx, `
		INSERT INTO task (
			id, job_id, parent_id, mutation, variables, worker_id, status, created,
			updated, started, finished, progress_current, progress_total, message
		)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?
		)
	`, row.ID, row.JobID, row.ParentID, row.Mutation, row.Variables, row.WorkerID, row.Status, row.Created,
		row.Updated, row.Started, row.Finished, row.ProgressCurrent, row.ProgressTotal, row.Message,
	); err != nil {
		return nil, err
	}
	return &TaskResolver{
		Q:       r,
		TaskRow: row,
	}, nil
}

func (r *QueryResolver) taskByID(ctx context.Context, id *string) (*TaskResolver, error) {
	t := &TaskResolver{}
	err := r.getRowByID(ctx, &t.TaskRow, `
		SELECT
			id, job_id, parent_id, mutation, variables, worker_id, status, created,
			updated, started, finished, progress_current, progress_total, message
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
			id, job_id, parent_id, mutation, variables, worker_id, status, created,
			updated, started, finished, progress_current, progress_total, message
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
		// No localization, fallback to mutation name.
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
