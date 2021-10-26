package main

import (
	"embed"
	"strings"

	"github.com/deref/exo/internal/about"
)

// SEE NOTE: [ABOUT_EMBEDS].

//go:embed VERSION
var Version string

//go:embed AMPLITUDE_API_KEY
var AmplitudeAPIKey string

//go:embed NOTICES.md
//go:embed LICENSE
//go:embed doc/licenses/ofl
var Notices embed.FS

func init() {
	about.Notices = Notices
	about.Version = strings.TrimSpace(Version)
	about.AmplitudeAPIKey = strings.TrimSpace(AmplitudeAPIKey)
}
