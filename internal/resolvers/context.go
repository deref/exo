package resolvers

import "context"

type contextKey int

const (
	taskTrackerKey contextKey = iota + 1
)

func ContextWithTaskTracker(ctx context.Context, tracker *TaskTracker) context.Context {
	return context.WithValue(ctx, taskTrackerKey, tracker)
}

func CurrentTaskTracker(ctx context.Context) *TaskTracker {
	tracker, _ := ctx.Value(taskTrackerKey).(*TaskTracker)
	return tracker
}
