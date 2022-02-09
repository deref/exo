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
	. "github.com/deref/exo/internal/scalars"
)

type TaskResolver struct {
	Q *RootResolver
	TaskRow
}

type TaskRow struct {
	ID              string     `db:"id"`
	JobID           string     `db:"job_id"`
	ParentID        *string    `db:"parent_id"`
	Mutation        string     `db:"mutation"`
	Arguments       JSONObject `db:"arguments"`
	WorkerID        *string    `db:"worker_id"`
	Status          string     `db:"status"`
	Created         Instant    `db:"created"`
	Updated         Instant    `db:"updated"`
	Started         *Instant   `db:"started"`
	Canceled        *Instant   `db:"canceled"`
	Finished        *Instant   `db:"finished"`
	ProgressCurrent *int32     `db:"progress_current"`
	ProgressTotal   *int32     `db:"progress_total"`
	ErrorMessage    *string    `db:"error_message"`
}

func (r *MutationResolver) CreateTask(ctx context.Context, args struct {
	ParentID  *string
	Mutation  string
	Arguments JSONObject
}) (*TaskResolver, error) {
	id := newTaskID()
	return r.createTask(ctx, id, args.ParentID, args.Mutation, args.Arguments)
}

var newTaskID = gensym.RandomBase32

func (r *MutationResolver) createJob(ctx context.Context, id string, mutation string, arguments map[string]interface{}) (*TaskResolver, error) {
	parentID := (*string)(nil)
	return r.createTask(ctx, id, parentID, mutation, arguments)
}

// The id is passed as a parameter to allow callers to use a pre-allocated id
// in a database field to establish a mutual exclusion lock.
func (r *MutationResolver) createTask(ctx context.Context, id string, parentID *string, mutation string, arguments map[string]interface{}) (*TaskResolver, error) {
	if id == "" {
		id = gensym.RandomBase32()
	}
	now := Now(ctx)
	row := TaskRow{
		ID:        id,
		ParentID:  parentID,
		Mutation:  mutation,
		Arguments: arguments,
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
	task, err := r.taskByID(ctx, &args.ID)
	if task == nil || err != nil {
		return task, err
	}
	if _, err := r.createEvent(ctx, task, "TaskStarted", ""); err != nil {
		return nil, fmt.Errorf("creating started event: %w", err)
	}
	return task, nil
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
	taskID := args.ID

	var errMessage *string
	if args.Error != nil {
		errMessage = args.Error
	}

	// XXX instead of erroring, await children to finish.
	// XXX cancel tree if parent failed before children.
	if errMessage == nil {
		outstanding := 0
		if err := r.DB.GetContext(ctx, &outstanding, `
			SELECT COUNT(id)
			FROM task
			WHERE parent_id = ?
			AND finished IS NULL
		`, args.ID,
		); err != nil {
			return nil, fmt.Errorf("counting outstanding tasks: %w", err)
		}
		if outstanding > 0 {
			msg := "task finished before children"
			errMessage = &msg
		}
	}

	now := Now(ctx)
	var status string
	if errMessage == nil {
		status = api.TaskStatusSuccess
	} else {
		status = api.TaskStatusFailure
	}
	var row TaskRow
	if err := r.DB.GetContext(ctx, &row, `
		UPDATE task
		SET updated = ?, finished = ?, status = ?, error_message = ?
		WHERE id = ?
		RETURNING *
	`, now, now, status, errMessage, taskID,
	); err != nil {
		return nil, fmt.Errorf("marking task as finished: %w", err)
	}
	task := &TaskResolver{
		Q:       r,
		TaskRow: row,
	}
	if _, err := r.createEvent(ctx, task, "TaskFinished", ""); err != nil {
		return nil, fmt.Errorf("creating finish event: %w", err)
	}
	return nil, nil
}

func (r *QueryResolver) TaskByID(ctx context.Context, args struct {
	ID string
}) (*TaskResolver, error) {
	return r.taskByID(ctx, &args.ID)
}

func (r *QueryResolver) taskByID(ctx context.Context, id *string) (*TaskResolver, error) {
	t := &TaskResolver{
		Q: r,
	}
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

func (r *QueryResolver) tasksByParentID(ctx context.Context, parentID string) ([]*TaskResolver, error) {
	var rows []TaskRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
		FROM task
		WHERE parent_id = ?
		ORDER BY task.id ASC
	`, parentID)
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

func (r *TaskResolver) Job() *JobResolver {
	return r.Q.jobByID(&r.JobID)
}

func (r *TaskResolver) Parent(ctx context.Context) (*TaskResolver, error) {
	return r.Q.taskByID(ctx, r.ParentID)
}

func (r *TaskResolver) Children(ctx context.Context) ([]*TaskResolver, error) {
	return r.Q.tasksByParentID(ctx, r.ID)
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

func (r *TaskResolver) Stream() *StreamResolver {
	return r.Q.streamForSource("Task", r.ID)
}

func (r *TaskResolver) eventPrototype(ctx context.Context) (row EventRow, err error) {
	// XXX set row's WorkspaceID, StackID, and ComponentID appropriately.
	row.JobID = &r.JobID
	row.TaskID = &r.ID
	return row, nil
}

func (r *TaskResolver) Message(ctx context.Context) (string, error) {
	if r.ErrorMessage != nil {
		return *r.ErrorMessage, nil
	}
	return r.Stream().Message(ctx)
}
