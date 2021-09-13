package logging

import (
	"context"
)

type contextKey int

const loggerKey contextKey = 1

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func CurrentLogger(ctx context.Context) Logger {
	return ctx.Value(loggerKey).(Logger)
}
