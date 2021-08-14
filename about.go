package exo

import (
	"embed"
)

const (
	CheckVersionEndpoint = "https://exo.deref.io/latest-version"
	UpdateScriptEndpoint = "https://exo.deref.io/install"
	TelemetryEndpoint    = "https://exo.deref.io/api/telemetry"
)

//go:embed VERSION
var Version string

//go:embed NOTICES.md
//go:embed LICENSE
//go:embed doc/licenses/ofl
var Notices embed.FS
