package main

import (
	"embed"

	"github.com/deref/exo/internal/about"
)

// SEE NOTE: [ABOUT_EMBEDS].

//go:embed VERSION
var Version string

//go:embed NOTICES.md
//go:embed LICENSE
//go:embed doc/licenses/ofl
var Notices embed.FS

func init() {
	about.Version = Version
	about.Notices = Notices
}
