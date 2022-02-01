package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVar(&applyFlags.Format, "format", "", "exo, compose, procfile")
}

var applyFlags struct {
	Format string
}

var applyCmd = &cobra.Command{
	Use:   "apply [flags] [manifest-file]",
	Short: "Applies a manifest in the current workspace",
	Long: `Applies a manifest in the current workspace.
	
If no manifest file is specified, a search is conducted in the current
directory in the following order of format preference:

	1. exo (Cue)
	2. exo (Legacy: HCL)
	3. compose
	4. procfile

The default Exo filename is 'exo.cue' and legacy Exo files have filename
'exo.hcl'.

Compose files may have one of the following names in order of preference:

	compose.yaml
	compose.yml
	docker-compose.yaml
	docker-compose.yml

The expected procfile name 'Procfile'.

If a manifest format will be guessed from the manifest filename.  This can be
overidden explicitly with the --format flag.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cl := newClient()
		kernel := cl.Kernel()
		workspace := requireCurrentWorkspace(ctx, cl)

		return apply(ctx, kernel, workspace, args)
	},
}

func apply(ctx context.Context, kernel api.Kernel, workspace api.Workspace, args []string) error {
	input := &api.ApplyInput{
		Format: applyFlags.Format,
	}
	if len(args) > 0 {
		manifestPath := args[0]
		input.ManifestPath = &manifestPath

		// We're not necessarily in the workspace root here,
		// so send the file contents too.
		bs, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return fmt.Errorf("reading manifest file: %w", err)
		}
		s := string(bs)
		input.Manifest = &s
	}

	output, err := workspace.Apply(ctx, input)
	if output != nil {
		for _, warning := range output.Warnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
		}
	}
	if err != nil {
		return err
	}
	return watchJob(ctx, output.JobID)
}
