package main

import (
	"os"

	"github.com/deref/exo/import/compose"
	"github.com/deref/exo/manifest"
	"github.com/deref/exo/util/cmdutil"
)

func main() {
	cfg, err := compose.Import(os.Stdin)
	if err != nil {
		cmdutil.Fatal(err)
	}
	if err := manifest.Generate(os.Stdout, cfg); err != nil {
		cmdutil.Fatal(err)
	}
}
