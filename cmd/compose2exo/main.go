package main

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/util/cmdutil"
)

func main() {
	importer := &compose.Loader{ProjectName: "unnamed"}
	res := importer.Load(os.Stdin)
	for _, warning := range res.Warnings {
		fmt.Fprintln(os.Stderr, warning)
	}
	if res.Err != nil {
		cmdutil.Fatal(res.Err)
	}
	if err := manifest.Generate(os.Stdout, res.Manifest); err != nil {
		cmdutil.Fatal(err)
	}
}
