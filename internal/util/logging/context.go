package logging

import (
	"context"

	"github.com/deref/exo/internal/util/contextutil"
)

type contextKey int

const loggerKey contextKey = 1

func init() {
	contextutil.RegisterContextCloner(func(src, dest context.Context) context.Context {
		return ContextWithLogger(dest, CurrentLogger(src))
	})
}

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func CurrentLogger(ctx context.Context) Logger {
	return ctx.Value(loggerKey).(Logger)
}
