package telemetry

type noOpTelemetry struct {
}

func (t *noOpTelemetry) IsEnabled() bool {
	return false
}

func (t *noOpTelemetry) LatestVersion() (string, error) {
	return "", nil
}

func (t *noOpTelemetry) StartSession() {
	// Do nothing.
}

func (t *noOpTelemetry) SendEvent(_ event) {
	// Do nothing.
}
