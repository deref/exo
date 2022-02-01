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

func Infof(ctx context.Context, format string, v ...interface{}) {
	logger := CurrentLogger(ctx)
	logger.Infof(format, v...)
}
