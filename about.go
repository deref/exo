package exo

import (
	"embed"
)

const (
	CheckVersionEndpoint = "https://exo.deref.io/latest-version"
	UpdateScriptEndpoint = "https://exo.deref.io/install"
)

//go:embed VERSION
var Version string

//go:embed AMPLITUDE_API_KEY
var AmplitudeAPIKey string

//go:embed NOTICES.md
//go:embed LICENSE
//go:embed doc/licenses/ofl
var Notices embed.FS
