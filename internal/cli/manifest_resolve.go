package cli

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestResolveCmd)
	manifestResolveCmd.Flags().StringVar(&manifestResolveFlags.Format, "format", "", "exo, compose, procfile")
}

var manifestResolveFlags struct {
	Format string
}

var manifestResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolves the workspace's manifest file.",
	Long:  `Resolves the workspace's manifest file and prints its path.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return nil
		}
		ctx := newContext()
		checkOrEnsureServer()

		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.ResolveManifest(ctx, &api.ResolveManifestInput{
			Format: manifestResolveFlags.Format,
		})
		if err != nil {
			return err
		}
		if output.Path != "" {
			fmt.Println(output.Path)
		}
		return nil
	},
}
