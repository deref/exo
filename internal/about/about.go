package about

import "embed"

const (
	CheckVersionEndpoint = "https://exo.deref.io/latest-version"
	UpdateScriptEndpoint = "https://exo.deref.io/install"
	TelemetryEndpoint    = "https://exo.deref.io/api/telemetry"
)

// NOTE [ABOUT_EMBEDS]: These come form go:embed tags, but the source files
// are in the repository root, and so must be set from there somehow.
// TODO: Untangle this.
var Version string
var Notices embed.FS
var AmplitudeAPIKey string
