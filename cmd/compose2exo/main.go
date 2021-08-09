package main

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/external/compose"
	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/cmdutil"
)

func main() {
	res := compose.Import(os.Stdin)
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
