package core

import _ "embed"

const (
	CheckVersionEndpoint = "https://download-page.deref.workers.dev/version.txt"
	UpdateScriptEndpoint = "https://download-page.deref.workers.dev/install.sh"
)

//go:embed VERSION
var Version string
