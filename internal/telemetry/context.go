package telemetry

import (
	"context"

	"github.com/deref/exo/internal/util/contextutil"
)

type contextKey int

const telemetryKey contextKey = 1

func init() {
	contextutil.RegisterContextCloner(func(src, dest context.Context) context.Context {
		return ContextWithTelemetry(dest, FromContext(src))
	})
}

func ContextWithTelemetry(ctx context.Context, telemetry Telemetry) context.Context {
	return context.WithValue(ctx, telemetryKey, telemetry)
}

func FromContext(ctx context.Context) Telemetry {
	return ctx.Value(telemetryKey).(Telemetry)
}
