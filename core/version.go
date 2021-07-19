package core

import _ "embed"

const (
	CheckVersionEndpoint = "https://exo.deref.io/latest-version"
	UpdateScriptEndpoint = "https://exo.deref.io/install"
)

//go:embed VERSION
var Version string
