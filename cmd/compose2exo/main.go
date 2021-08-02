package main

import (
	"fmt"
	"os"

	"github.com/deref/exo/import/compose"
	"github.com/deref/exo/manifest"
	"github.com/deref/exo/util/cmdutil"
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
