package cli

import (
	"fmt"

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
		ctx := cmd.Context()
		var q struct {
			Workspace *struct {
				Manifest *struct {
					HostPath string
				}
			} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
		}

		vars := map[string]interface{}{}
		if manifestResolveFlags.Format != "" {
			vars["format"] = manifestResolveFlags.Format
		}
		mustQueryWorkspace(ctx, &q, vars)
		manifest := q.Workspace.Manifest
		if manifest != nil {
			fmt.Println(manifest.HostPath)
		}
		return nil
	},
}
