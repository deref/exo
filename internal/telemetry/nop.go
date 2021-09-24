package telemetry

import "context"

type Nop struct {
}

func (t *Nop) IsEnabled() bool {
	return false
}

func (t *Nop) LatestVersion(_ context.Context) (string, error) {
	return "", nil
}

func (t *Nop) StartSession(_ context.Context) {
	// Do nothing.
}

func (t *Nop) SendEvent(_ context.Context, _ event) {
	// Do nothing.
}

func (t *Nop) RecordOperation(_ OperationInvocation) {
	// Do nothing.
}
