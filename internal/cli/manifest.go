package cli

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.AddCommand(makeHelpSubcmd())
}

var manifestCmd = &cobra.Command{
	Use:    "manifest",
	Short:  "Manifest tools",
	Long:   `Contains subcommands for working with manifests`,
	Hidden: true,
	Args:   cobra.NoArgs,
}

func loadManifest(name string) (*manifest.Manifest, error) {
	// TODO: Support other formats here too.
	loader := &manifest.Loader{}
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}
	defer f.Close()
	res := loader.Load(f)
	for _, warning := range res.Warnings {
		cmdutil.Warnf("%s\n", warning)
	}
	return res.Manifest, res.Err
}
