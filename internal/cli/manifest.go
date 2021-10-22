package cli

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/util/term"
	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.AddCommand(makeHelpSubcmd())
	manifestCmd.PersistentFlags().StringVar(&manifestFlags.Format, "format", "", "exo, compose, procfile")
}

var manifestFlags struct {
	Format string
}

var manifestCmd = &cobra.Command{
	Use:    "manifest",
	Short:  "Manifest tools",
	Long:   `Contains subcommands for working with manifests`,
	Hidden: true,
	Args:   cobra.NoArgs,
}

func loadManifest(name string) (*exohcl.Manifest, error) {
	bs, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}

	loader := &manifest.Loader{
		WorkspaceName: "unnamed",
		Format:        manifestFlags.Format,
		Filename:      name,
		Bytes:         bs,
	}
	return loader.Load()
}

func writeManifestError(w io.Writer, err error) error {
	var diags hcl.Diagnostics
	if !errors.As(err, &diags) {
		return err
	}
	files := map[string]*hcl.File{} // TODO: Populate map for .hcl input files.
	width, _ := term.GetSize()
	enableColor := true // https://github.com/deref/exo/issues/179
	diagWr := hcl.NewDiagnosticTextWriter(w, files, uint(width), enableColor)
	return diagWr.WriteDiagnostics(diags)
}
