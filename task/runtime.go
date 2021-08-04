package task

import (
	"context"

	"github.com/deref/exo/task/api"
	"golang.org/x/sync/errgroup"
)

type TaskTracker struct {
	Store api.TaskStore
}

type Task struct {
	ctx context.Context
	tt  *TaskTracker
	id  string
	def Definition
	eg  errgroup.Group
}

func (tt *TaskTracker) CreateTask(ctx context.Context, def Definition) *Task {
	var parent *Task
	return tt.createTask(ctx, parent, def)
}

func (tt *TaskTracker) CreateTaskFunc(ctx context.Context, name string, run func(*Task) error) *Task {
	var parent *Task
	return tt.createTask(ctx, parent, NewTaskFunc(name, run))
}

func (tt *TaskTracker) createTask(ctx context.Context, parent *Task, def Definition) *Task {
	input := &api.CreateTaskInput{
		Name: def.Name(),
	}
	if parent != nil {
		input.ParentID = &parent.id
	}
	output, err := tt.Store.CreateTask(ctx, input)
	if err != nil {
		panic(err) // XXX
	}
	return &Task{
		tt:  tt,
		id:  output.ID,
		def: def,
	}
}

func (tt *TaskTracker) RunTask(ctx context.Context, def Definition) error {
	t := tt.CreateTask(ctx, def)
	return t.Run()
}

func (tt *TaskTracker) RunFunc(ctx context.Context, name string, f func(*Task) error) error {
	t := tt.CreateTaskFunc(ctx, name, f)
	return t.Run()
}

func (t *Task) Context() context.Context {
	return t.ctx
}

func (t *Task) reportStatus(status string) {
	_, err := t.tt.Store.UpdateTask(t.ctx, &api.UpdateTaskInput{
		Status: &status,
	})
	if err != nil {
		panic(err) // XXX
	}
}

func (t *Task) Run() (err error) {
	t.reportStatus(api.StatusRunning)
	defer func() {
		if err == nil {
			t.reportStatus(api.StatusFailure)
		} else {
			t.reportStatus(api.StatusSuccess)
		}
	}()
	err = t.def.RunTask(t)
	err2 := t.eg.Wait()
	if err == nil {
		err = err2
	}
	return
}

func (t *Task) Go(def Definition) {
	subtask := t.tt.createTask(t.ctx, t, def)
	t.eg.Go(func() error {
		return subtask.Run()
	})
}

func (t *Task) GoFunc(name string, f func(*Task) error) {
	t.Go(NewTaskFunc(name, f))
}
