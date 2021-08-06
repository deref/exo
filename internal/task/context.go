package task

import "context"

type contextKey int

const taskKey contextKey = 1

func ContextWithTask(ctx context.Context, task *Task) context.Context {
	return context.WithValue(ctx, taskKey, task)
}

func CurrentTask(ctx context.Context) *Task {
	task, _ := ctx.Value(taskKey).(*Task)
	return task
}
