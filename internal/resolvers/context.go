package resolvers

import "context"

type contextKey int

const (
	taskIDKey contextKey = iota + 1
)

type TaskContext struct {
	JobID    string
	ID       string
	WorkerID string
}

func ContextWithTask(ctx context.Context, tc TaskContext) context.Context {
	return context.WithValue(ctx, taskIDKey, &tc)
}

func CurrentTask(ctx context.Context) *TaskContext {
	tc, _ := ctx.Value(taskIDKey).(*TaskContext)
	return tc
}
