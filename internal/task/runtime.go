package task

import (
	"context"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/util/logging"
	"golang.org/x/sync/errgroup"
)

type TaskTracker struct {
	Store  api.TaskStore
	Logger logging.Logger
}

type Task struct {
	context.Context
	cancel       func()
	tt           *TaskTracker
	jobID        string
	id           string
	status       string
	eg           errgroup.Group
	err          error
	reportingErr error
}

// ID of this task.  Might be empty string if task tracking failed.
func (t *Task) ID() string {
	return t.id
}

// Job ID associated with this task. Might be empty string if task tracking failed.
func (t *Task) JobID() string {
	return t.jobID
}

func (t *Task) Err() error {
	return t.err
}

// Creates a task in the Pending state. Call task.Start() on it, or use
// StartTask method for convenience.
func (tt *TaskTracker) CreateTask(ctx context.Context, name string) *Task {
	var parent *Task
	return tt.createTask(ctx, parent, name)
}

func (tt *TaskTracker) createTask(ctx context.Context, parent *Task, name string) *Task {
	ctx, cancel := context.WithCancel(ctx)
	input := &api.CreateTaskInput{
		Name: name,
	}
	if parent != nil {
		input.ParentID = &parent.id
	}
	t := &Task{
		cancel: cancel,
		tt:     tt,
		status: api.StatusPending,
	}
	t.Context = ContextWithTask(ctx, t)
	output, err := tt.Store.CreateTask(ctx, input)
	if err != nil {
		t.reportingErr = err
	} else {
		t.id = output.ID
		t.jobID = output.JobID
	}
	return t
}

// Creates and starts a task. See `Task.Start()` for usage instructions.
func (tt *TaskTracker) StartTask(ctx context.Context, name string) *Task {
	var parent *Task
	return tt.startTask(ctx, parent, name)
}

func (tt *TaskTracker) startTask(ctx context.Context, parent *Task, name string) *Task {
	task := tt.createTask(ctx, parent, name)
	task.Start()
	return task
}

// Marks a task as started. The caller is responsible for also calling .Finish().
func (t *Task) Start() {
	if t.status != api.StatusPending {
		panic("cannot start task that is not in pending state")
	}
	t.updateTask(api.StatusRunning, "", "", 0, 0)
}

// Marks this Task as failed, which will be reported by Finish().
func (t *Task) Fail(err error) {
	if t.err == nil {
		t.err = err
	}
}

// Waits for any subtasks to complete, then reports this task as complete.
// Status will be failed if any subtask fails, or if Fail() is called.
// It is safe to call Finish() multiple times.
func (t *Task) Finish() error {
	cancel := t.cancel
	if cancel == nil {
		return t.err
	}
	defer cancel()
	t.cancel = nil
	var err error
	if t.err != nil {
		err = t.err
		cancel()
	} else {
		err = t.eg.Wait()
	}
	status := api.StatusSuccess
	message := ""
	if err != nil {
		status = api.StatusFailure
		message = err.Error()
	}
	finished := chrono.NowString(t)
	t.updateTask(status, message, finished, 0, 0)
	if t.reportingErr != nil {
		t.tt.Logger.Infof("error reporting task progress: %v", t.reportingErr)
	}
	return err
}

func (t *Task) updateTask(status string, message string, finished string, current, total int) {
	if t.id == "" {
		return
	}
	input := api.UpdateTaskInput{
		ID:     t.id,
		Status: &status,
	}
	if message != "" {
		input.Status = &message
	}
	if finished != "" {
		input.Finished = &finished
	}
	if total > 0 {
		input.Progress = &api.TaskProgress{
			Current: current,
			Total:   total,
		}
	}
	if _, err := t.tt.Store.UpdateTask(t, &input); err != nil {
		if t.reportingErr == nil {
			t.reportingErr = err
		}
	}
}

func (t *Task) ReportMessage(message string) {
	t.updateTask("", message, "", 0, 0)
}

func (t *Task) ReportProgress(current, total int) {
	t.updateTask("", "", "", current, total)
}

func (t *Task) Wait() error {
	return t.eg.Wait()
}

func (t *Task) StartChild(name string) *Task {
	return t.tt.startTask(t, t, name)
}

func (t *Task) RunChild(name string, f func(task *Task) error) error {
	task := t.StartChild(name)
	if err := f(task); err != nil {
		task.Fail(err)
	}
	return task.Finish()
}

func (t *Task) Go(name string, f func(task *Task) error) {
	t.eg.Go(func() error {
		return t.RunChild(name, f)
	})
}
