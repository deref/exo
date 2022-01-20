package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/chrono"
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

func (r *MutationResolver) NewTask(ctx context.Context, args struct {
	ParentID  *string
	Mutation  string
	Variables string
}) (*TaskResolver, error) {
	id := newTaskID()
	return r.newTask(ctx, id, args.ParentID, args.Mutation, args.Variables)
}

var newTaskID = gensym.RandomBase32

// The id is passed as a parameter to allow callers to use a pre-allocated id
// in a database field to establish a mutual exclusion lock.
func (r *MutationResolver) newTask(ctx context.Context, id string, parentID *string, mutation string, variables string) (*TaskResolver, error) {
	if id == "" {
		panic("id is required")
	}
	now := Now(ctx)
	row := TaskRow{
		ID:        id,
		ParentID:  parentID,
		Mutation:  mutation,
		Variables: variables,
		Status:    api.TaskStatusPending,
		Created:   now,
		Updated:   now,
	}
	if parentID == nil {
		row.JobID = id
	} else {
		parent, err := r.taskByID(ctx, parentID)
		if err != nil {
			return nil, fmt.Errorf("resolving parent: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("no such parent: %q", *parentID)
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

func (r *MutationResolver) AcquireTask(ctx context.Context, args struct {
	WorkerID      string
	TimeoutMillis *int
}) (*TaskResolver, error) {
	if args.TimeoutMillis != nil {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*args.TimeoutMillis)*time.Microsecond)
		defer cancel()
	}
	var row TaskRow
	delay := 1
	for {
		fmt.Println("querying")
		err := r.DB.GetContext(ctx, &row, `
		UPDATE task
		SET worker_id = ?
		WHERE id IN (
			SELECT id
			FROM task
			WHERE worker_id IS NULL
		)
		RETURNING *
	`, args.WorkerID)
		if errors.Is(err, sql.ErrNoRows) {
			chrono.Sleep(ctx, time.Duration(delay)*time.Millisecond)
			delay *= 2
			if delay > 1000 {
				delay = 1000
			}
			continue
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		break
	}
	return &TaskResolver{
		Q:       r,
		TaskRow: row,
	}, nil
}

func (r *MutationResolver) StartTask(ctx context.Context, args struct {
	ID       string
	WorkerID string
}) (*TaskResolver, error) {
	now := Now(ctx)
	res, err := r.DB.ExecContext(ctx, `
		UPDATE task
		SET worker_id = ?, started = ?
		WHERE id = ?
		AND (worker_id = ? OR worker_id IS NULL)
	`, args.WorkerID, now, args.ID, args.WorkerID)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if n != 1 {
		return nil, errors.New("task not available")
	}
	return r.taskByID(ctx, &args.ID)
}

func (r *MutationResolver) UpdateTask(ctx context.Context, args struct {
	ID       string
	Message  *string
	Progress *ProgressInput
}) (*VoidResolver, error) {
	return r.updateTask(ctx, args.ID, args.Message, args.Progress)
}

func (r *MutationResolver) updateTask(ctx context.Context, id string, message *string, progress *ProgressInput) (*VoidResolver, error) {
	now := Now(ctx)
	var progressCurrent, progressTotal *int32
	if progress != nil {
		progressCurrent = &progress.Current
		progressTotal = &progress.Total
	}
	_, err := r.DB.ExecContext(ctx, `
		UPDATE task
		SET
			updated = ?,
			message = COALESCE(?, message),
			progress_current = COALESCE(?, progress_current),
			progress_total = COALESCE(?, progress_total)
		WHERE id = ?
	`, now, message, progressCurrent, progressTotal)
	return nil, err
}

func (r *MutationResolver) FinishTask(ctx context.Context, args struct {
	ID    string
	Error *string
}) (*VoidResolver, error) {
	now := Now(ctx)
	var status string
	if args.Error == nil {
		status = api.TaskStatusSuccess
	} else {
		status = api.TaskStatusFailure
	}
	_, err := r.DB.ExecContext(ctx, `
		UPDATE task
		SET updated = ?, finished = ?, status = ?, message = ?
		WHERE id = ?
	`, now, now, status, args.Error, args.ID)
	return nil, err
}

func (r *QueryResolver) taskByID(ctx context.Context, id *string) (*TaskResolver, error) {
	t := &TaskResolver{}
	err := r.getRowByKey(ctx, &t.TaskRow, `
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

func (r *MutationResolver) cancelTask(ctx context.Context, id string) error {
	r.Logger.Infof("TODO: Implement cancelTask")
	return nil
}
