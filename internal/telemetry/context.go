package telemetry

import "context"

type contextKey int

const telemetryKey contextKey = 1

func ContextWithTelemetry(ctx context.Context, telemetry Telemetry) context.Context {
	return context.WithValue(ctx, telemetryKey, telemetry)
}

func FromContext(ctx context.Context) Telemetry {
	return ctx.Value(telemetryKey).(Telemetry)
}
