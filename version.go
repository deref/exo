package exo

import _ "embed"

const (
	CheckVersionEndpoint = "https://exo.deref.io/latest-version"
	UpdateScriptEndpoint = "https://exo.deref.io/install"
	TelemetryEndpoint    = "https://exo.deref.io/api/telemetry"
)

//go:embed VERSION
var Version string
