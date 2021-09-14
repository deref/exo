package telemetry

import "context"

type noOpTelemetry struct {
}

func (t *noOpTelemetry) IsEnabled() bool {
	return false
}

func (t *noOpTelemetry) LatestVersion(_ context.Context) (string, error) {
	return "", nil
}

func (t *noOpTelemetry) StartSession(_ context.Context) {
	// Do nothing.
}

func (t *noOpTelemetry) SendEvent(_ context.Context, _ Event) {
	// Do nothing.
}

func (t *noOpTelemetry) RecordOperation(_ OperationInvocation) {
	// Do nothing.
}
