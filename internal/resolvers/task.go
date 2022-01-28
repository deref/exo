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
	"github.com/deref/exo/internal/util/jsonutil"
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
	Arguments       string   `db:"arguments"`
	WorkerID        *string  `db:"worker_id"`
	Status          string   `db:"status"`
	Created         Instant  `db:"created"`
	Updated         Instant  `db:"updated"`
	Started         *Instant `db:"started"`
	Canceled        *Instant `db:"canceled"`
	Finished        *Instant `db:"finished"`
	ProgressCurrent *int32   `db:"progress_current"`
	ProgressTotal   *int32   `db:"progress_total"`
}

func (r *MutationResolver) CreateTask(ctx context.Context, args struct {
	ParentID  *string
	Mutation  string
	Arguments string
}) (*TaskResolver, error) {
	id := newTaskID()
	var taskArgs map[string]interface{}
	if err := jsonutil.UnmarshalString(args.Arguments, &taskArgs); err != nil {
		return nil, fmt.Errorf("unmarshaling task arguments: %w", err)
	}
	return r.createTask(ctx, id, args.ParentID, args.Mutation, taskArgs)
}

var newTaskID = gensym.RandomBase32

func (r *MutationResolver) createJob(ctx context.Context, id string, mutation string, args map[string]interface{}) (*TaskResolver, error) {
	parentID := (*string)(nil)
	return r.createTask(ctx, id, parentID, mutation, args)
}

// The id is passed as a parameter to allow callers to use a pre-allocated id
// in a database field to establish a mutual exclusion lock.
func (r *MutationResolver) createTask(ctx context.Context, id string, parentID *string, mutation string, args map[string]interface{}) (*TaskResolver, error) {
	if id == "" {
		panic("id is required")
	}
	now := Now(ctx)
	row := TaskRow{
		ID:        id,
		ParentID:  parentID,
		Mutation:  mutation,
		Arguments: jsonutil.MustMarshalString(args),
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
	if err := r.insertRow(ctx, "task", row); err != nil {
		return nil, err
	}
	return &TaskResolver{
		Q:       r,
		TaskRow: row,
	}, nil
}

func (r *MutationResolver) AcquireTask(ctx context.Context, args struct {
	WorkerID string
	JobID    *string
	Timeout  *int32
}) (*TaskResolver, error) {
	if args.Timeout != nil {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*args.Timeout)*time.Microsecond)
		defer cancel()
	}
	var row TaskRow
	delay := 1
	for {
		err := r.DB.GetContext(ctx, &row, `
		UPDATE task
		SET worker_id = ?
		WHERE id IN (
			SELECT id
			FROM task
			WHERE worker_id IS NULL
			AND COALESCE(?, job_id) = job_id
		)
		RETURNING *
	`, args.WorkerID, args.JobID)
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
	WorkerID string
	Progress *ProgressInput
}) (*TaskResolver, error) {
	return r.updateTask(ctx, args.ID, args.WorkerID, args.Progress)
}

func (r *MutationResolver) updateTask(ctx context.Context, id string, workerID string, progress *ProgressInput) (*TaskResolver, error) {
	now := Now(ctx)
	var progressCurrent, progressTotal *int32
	if progress != nil {
		progressCurrent = &progress.Current
		progressTotal = &progress.Total
	}
	var row TaskRow
	err := r.DB.GetContext(ctx, &row, `
		UPDATE task
		SET
			updated = ?,
			progress_current = COALESCE(?, progress_current),
			progress_total = COALESCE(?, progress_total)
		WHERE id = ?
		AND worker_id = ?
		RETURNING *
	`, now, progressCurrent, progressTotal, id, workerID)
	if err != nil {
		return nil, err
	}
	return &TaskResolver{
		Q:       r,
		TaskRow: row,
	}, nil
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
		SELECT *
		FROM task
		WHERE id = ?
	`, id)
	if t.ID == "" {
		t = nil
	}
	return t, err
}

func (r *QueryResolver) jobByID(ctx context.Context, id *string) (*TaskResolver, error) {
	task, err := r.taskByID(ctx, id)
	if task != nil && task.JobID != task.ID {
		task = nil
	}
	return task, err
}

func (r *QueryResolver) AllTasks(ctx context.Context) ([]*TaskResolver, error) {
	var rows []TaskRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
		FROM task
		ORDER BY task.id ASC
	`)
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

func (r *QueryResolver) TasksByJobID(ctx context.Context, args struct {
	JobID string
}) ([]*TaskResolver, error) {
	return r.tasksByJobID(ctx, args.JobID)
}

func (r *QueryResolver) tasksByJobID(ctx context.Context, jobID string) ([]*TaskResolver, error) {
	return r.tasksByJobIDs(ctx, []string{jobID})
}

func (r *QueryResolver) TasksByJobIDs(ctx context.Context, args struct {
	JobIDs []string
}) ([]*TaskResolver, error) {
	return r.tasksByJobIDs(ctx, args.JobIDs)
}

func (r *QueryResolver) tasksByJobIDs(ctx context.Context, jobIDs []string) ([]*TaskResolver, error) {
	var rows []TaskRow
	if len(jobIDs) > 0 {
		query, args := mustSqlIn(`
			SELECT *
			FROM task
			WHERE job_id IN (?)
			ORDER BY task.id ASC
		`, jobIDs)
		err := r.DB.SelectContext(ctx, &rows, query, args...)
		if err != nil {
			return nil, err
		}
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

func (r *TaskResolver) Stream(ctx context.Context) (*StreamResolver, error) {
	return r.Q.streamByOwner(ctx, "Task", r.ID)
}

func (r *TaskResolver) Message(ctx context.Context) (string, error) {
	stream, err := r.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("resolving stream: %w", err)
	}
	return stream.Message(ctx)
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
