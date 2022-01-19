package resolvers

import "context"

type contextKey int

const (
	taskIDKey contextKey = iota + 1
)

func ContextWithTaskID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, taskIDKey, id)
}

func CurrentTaskID(ctx context.Context) *string {
	if taskID, ok := ctx.Value(taskIDKey).(string); ok {
		return &taskID
	}
	return nil
}
