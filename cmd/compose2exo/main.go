package main

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/util/cmdutil"
)

func main() {
	cmd, err := cmdutil.ParseArgs(os.Args)
	if err != nil {
		cmdutil.Fatalf("parsing arguments: %v", err)
	}
	projectName := "imported"
	if flagProjectName, ok := cmd.Flags["project-name"]; ok {
		projectName = flagProjectName
	}
	importer := compose.Importer{ProjectName: projectName}
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
