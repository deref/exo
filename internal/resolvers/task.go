package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/scalars"
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
	Key             *string    `db:"key"`
	WorkerID        *string    `db:"worker_id"`
	Created         Instant    `db:"created"`
	Updated         Instant    `db:"updated"`
	Started         *Instant   `db:"started"`
	Canceled        *Instant   `db:"canceled"`
	Finished        *Instant   `db:"finished"`
	Completed       *Instant   `db:"completed"`
	ProgressCurrent int32      `db:"progress_current"`
	ProgressTotal   int32      `db:"progress_total"`
	Error           *string    `db:"error"`
}

func (r *MutationResolver) CreateTask(ctx context.Context, args struct {
	Mutation  string
	Arguments JSONObject
	Key       *string
}) (*TaskResolver, error) {
	return r.createOrEnsureTask(ctx, args.Mutation, args.Arguments, args.Key)
}

func (r *MutationResolver) createRootTask(ctx context.Context, mutation string, arguments map[string]any) (*TaskResolver, error) {
	row := newTaskPrototype(ctx, mutation, arguments)
	row.JobID = row.ID
	if err := r.insertRow(ctx, "task", row); err != nil {
		return nil, err
	}
	return &TaskResolver{
		Q:       r,
		TaskRow: row,
	}, nil
}

func (r *MutationResolver) createTask(ctx context.Context, mutation string, arguments map[string]any) (*TaskResolver, error) {
	return r.createOrEnsureTask(ctx, mutation, arguments, nil)
}

func (r *MutationResolver) ensureTask(ctx context.Context, mutation string, arguments map[string]any, key string) error {
	_, err := r.createOrEnsureTask(ctx, mutation, arguments, &key)
	if isSqlConflict(err) {
		err = nil
	}
	return err
}

func (r *MutationResolver) createOrEnsureTask(ctx context.Context, mutation string, arguments map[string]any, key *string) (*TaskResolver, error) {
	ctxVars := api.CurrentContextVariables(ctx)
	if ctxVars == nil || ctxVars.TaskID == "" {
		return nil, errors.New("create task outside of job execution context")
	}
	row := newTaskPrototype(ctx, mutation, arguments)
	row.JobID = ctxVars.JobID
	row.ParentID = &ctxVars.TaskID
	if err := r.insertRow(ctx, "task", row); err != nil {
		return nil, err
	}
	return &TaskResolver{
		Q:       r,
		TaskRow: row,
	}, nil
}

type TaskInput struct {
	Mutation  string
	Arguments map[string]any
}

func (r *MutationResolver) createTasks(ctx context.Context, inputs []TaskInput) ([]*TaskResolver, error) {
	// TODO: Create tasks in bulk.
	tasks := make([]*TaskResolver, len(inputs))
	for i, input := range inputs {
		var err error
		tasks[i], err = r.createTask(ctx, input.Mutation, input.Arguments)
		if err != nil {
			return tasks, err
		}
	}
	return tasks, nil
}

func newTaskPrototype(ctx context.Context, mutation string, arguments map[string]any) TaskRow {
	now := Now(ctx)
	return TaskRow{
		ID:        gensym.RandomBase32(),
		Mutation:  mutation,
		Arguments: arguments,
		Created:   now,
		Updated:   now,
	}
}

func (r *MutationResolver) AcquireTask(ctx context.Context, args struct {
	WorkerID string
	JobID    *string
}) (*TaskResolver, error) {
	delay := 1
	for {
		// Attempt to assign a worker.
		var attemptedID string
		err := r.db.GetContext(ctx, &attemptedID, `
			UPDATE task
			SET worker_id = ?
			WHERE id IN (
				SELECT id
				FROM task
				WHERE worker_id IS NULL
				AND started IS NULL
				AND COALESCE(?, job_id) = job_id
				ORDER BY random()
				LIMIT 1
			)
			RETURNING id
		`, args.WorkerID, args.JobID)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		if err != nil {
			return nil, fmt.Errorf("assigning task worker: %w", err)
		}

		if attemptedID == "" {
			// There are no available tasks.

			// When scoped to just one job, and it is complete, there are no more
			// tasks we can possibly acquire, so return normally.
			if args.JobID != nil {
				completed, err := r.isJobCompleted(ctx, *args.JobID)
				if err != nil {
					return nil, fmt.Errorf("checking job completion: %w", err)
				}
				if completed {
					return nil, nil
				}
			}

			if err := chrono.Sleep(ctx, time.Duration(delay)*time.Millisecond); err != nil {
				return nil, err
			}
			delay *= 2
			if delay > 1000 {
				delay = 1000
			}
			continue
		}

		// Return the task, if we won the race to acquire it.
		{
			var row TaskRow
			err := r.db.GetContext(ctx, &row, `
				SELECT *
				FROM task
				WHERE id = ?
				AND worker_id = ?
				LIMIT 1
			`, attemptedID, args.WorkerID)
			if errors.Is(err, sql.ErrNoRows) {
				// We lost, another worker won. Try again.
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("checking task acquisition: %w", err)
			}

			return &TaskResolver{
				Q:       r,
				TaskRow: row,
			}, nil
		}
	}
}

func (r *MutationResolver) StartTask(ctx context.Context, args struct {
	ID       string
	WorkerID string
}) (*TaskResolver, error) {
	now := Now(ctx)
	res, err := r.db.ExecContext(ctx, `
		UPDATE task
		SET worker_id = ?, started = ?
		WHERE id = ?
		AND (worker_id = ? OR worker_id IS NULL)
	`, args.WorkerID, now, args.ID, args.WorkerID)
	if err != nil {
		return nil, err
	}
	if rowsAffected(res) != 1 {
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
	var progressCurrent, progressTotal int32
	if progress != nil {
		progressCurrent = progress.Current
		progressTotal = progress.Total
	}
	var row TaskRow
	err := r.db.GetContext(ctx, &row, `
		UPDATE task
		SET
			updated = ?,
			progress_current = MAX(?, progress_current),
			progress_total = MAX(?, progress_total)
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

	now := Now(ctx)
	var row TaskRow
	if err := r.db.GetContext(ctx, &row, `
		UPDATE task
		SET
			updated = ?,
			finished = ?,
			error = COALESCE(error, ?)
		WHERE id = ?
		AND finished IS NULL
		RETURNING *
	`, now, now, args.Error, taskID,
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
	if err := r.maybeCompleteTask(ctx, task); err != nil {
		return nil, fmt.Errorf("completing task: %w", err)
	}
	return nil, nil
}

// Walks up the task hierarchy, marking tasks as complete if all of their
// subtasks are complete.  When a subtask fails, propagates its error to the
// parent. If more than one task or its children fail, the first error is
// preserved.
func (r *MutationResolver) maybeCompleteTask(ctx context.Context, task *TaskResolver) error {
	now := Now(ctx)
	res, err := r.db.ExecContext(ctx, `
		UPDATE task
		SET
			completed = ?,
			error = COALESCE(error, (
				SELECT child.error
				FROM task AS child
				WHERE child.parent_id = task.id
				AND error IS NOT NULL
				LIMIT 1
			))
		WHERE id = ?
		AND completed IS NULL
		AND finished IS NOT NULL
		AND 0 == (
			SELECT count(child.id)
			FROM task AS child
			WHERE child.parent_id = task.id
			AND completed IS NULL
		)
	`, now, task.ID)
	if err != nil {
		return fmt.Errorf("marking task as complete: %w", err)
	}
	if rowsAffected(res) == 0 {
		return nil
	}
	if _, err := r.createEvent(ctx, task, "TaskCompleted", ""); err != nil {
		return fmt.Errorf("creating TaskCompleted event: %w", err)
	}
	if task.ParentID == nil {
		if _, err := r.createEvent(ctx, task.Job(), "JobCompleted", ""); err != nil {
			return fmt.Errorf("creating JobCompleted event: %w", err)
		}
		return nil
	}
	parent, err := task.Parent(ctx)
	if err != nil {
		return fmt.Errorf("resolving parent: %w", err)
	}
	return r.maybeCompleteTask(ctx, parent)
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

func taskRowsToResolvers(r *RootResolver, rows []TaskRow) []*TaskResolver {
	resolvers := make([]*TaskResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &TaskResolver{
			Q:       r,
			TaskRow: row,
		}
	}
	return resolvers
}

func (r *QueryResolver) AllTasks(ctx context.Context) ([]*TaskResolver, error) {
	var rows []TaskRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM task
		ORDER BY task.id ASC
	`)
	if err != nil {
		return nil, err
	}
	return taskRowsToResolvers(r, rows), nil
}

func (r *QueryResolver) tasksByParentID(ctx context.Context, parentID string) ([]*TaskResolver, error) {
	var rows []TaskRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM task
		WHERE parent_id = ?
		ORDER BY task.id ASC
	`, parentID)
	if err != nil {
		return nil, err
	}
	return taskRowsToResolvers(r, rows), nil
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
		err := r.db.SelectContext(ctx, &rows, query, args...)
		if err != nil {
			return nil, err
		}
	}
	return taskRowsToResolvers(r, rows), nil
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

func (r *TaskResolver) Label(ctx context.Context) (string, error) {
	methodName := fmt.Sprintf("%s%s_label", strings.ToUpper(r.Mutation[0:1]), r.Mutation[1:])
	method := reflect.ValueOf(r.Q).MethodByName(methodName)
	if !method.IsValid() {
		// No localization method, fallback to mutation name.
		return r.Mutation, nil
	}
	methodType, _ := reflect.TypeOf(r.Q).MethodByName(methodName)
	args := reflect.New(methodType.Type.In(2))
	if err := scalars.DecodeStruct((map[string]any)(r.Arguments), args.Interface()); err != nil {
		return "", fmt.Errorf("decoding arguments: %w", err)
	}
	// TODO: Use more flexible calling conventions, like other resolver methods.
	// 1) Optional context.
	// 2) Optional, strongly typed arguments.
	// 3) Optional error return value.
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), args.Elem()})
	label := res[0].Interface().(string)
	err, _ := res[1].Interface().(error)
	return label, err
}

func (r *TaskResolver) Progress() (*ProgressResolver, error) {
	if r.ProgressCurrent <= 0 || r.ProgressTotal <= 0 {
		return nil, nil
	}
	return &ProgressResolver{
		Current: r.ProgressCurrent,
		Total:   r.ProgressTotal,
	}, nil
}

func (r *TaskResolver) Stream() *StreamResolver {
	return r.Q.streamForSource("Task", r.ID)
}

func (r *TaskResolver) eventPrototype(ctx context.Context) (row EventRow, err error) {
	// XXX set row's WorkspaceID, StackID, and ComponentID appropriately.
	row.SourceType = "Task"
	row.JobID = &r.JobID
	row.TaskID = &r.ID
	return row, nil
}

func (r *TaskResolver) Message(ctx context.Context) (string, error) {
	if r.Error != nil {
		return fmt.Sprintf("error: %s", *r.Error), nil
	}
	return r.Stream().Message(ctx)
}

func (r *MutationResolver) CancelTask(ctx context.Context, args struct {
	ID string
}) (*VoidResolver, error) {
	err := r.cancelTask(ctx, args.ID)
	return nil, err
}

// See also cancelJob and cancelSubtasks.
func (r *MutationResolver) cancelTask(ctx context.Context, id string) error {
	now := Now(ctx)
	_, err := r.db.ExecContext(ctx, `
		UPDATE task
		SET canceled = COALESCE(canceled, ?)
		WHERE id IN (
			WITH RECURSIVE rec (id) AS (
				SELECT ?
				UNION
				SELECT id FROM task, rec WHERE task.parent_id = rec.id
			)
			SELECT id FROM rec
		) AND finished IS NULL
	`, now, id)
	return err
}

// See also cancelJob and cancelTask.
func (r *MutationResolver) cancelSubtasks(ctx context.Context, parentTaskID string) error {
	now := Now(ctx)
	_, err := r.db.ExecContext(ctx, `
		UPDATE task
		SET canceled = COALESCE(canceled, ?)
		WHERE id IN (
			WITH RECURSIVE rec (id) AS (
				SELECT id FROM task WHERE parent_id = ?
				UNION
				SELECT id FROM task, rec WHERE task.parent_id = rec.id
			)
			SELECT id FROM rec
		) AND finished IS NULL
	`, now, parentTaskID)
	return err
}

func (r *QueryResolver) isTaskCompleted(ctx context.Context, taskID string) (completed bool, err error) {
	err = r.db.GetContext(ctx, &completed, `
		SELECT completed IS NOT NULL
		FROM task
		WHERE id = ?
	`, taskID)
	return
}

func (r *TaskResolver) Successful() (*bool, error) {
	if r.Completed == nil {
		return nil, nil
	}
	ok := r.Error == nil
	return &ok, nil
}
