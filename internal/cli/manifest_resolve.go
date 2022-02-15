package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestResolveCmd)
	addManifestFormatFlag(manifestResolveCmd, &manifestResolveFlags.Format)
}

var manifestResolveFlags struct {
	Format string
}

var manifestResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Finds a manifest file.",
	Long: `Finds a manifest file in the current workspace and prints its path.

Ignores the workspaces's configured manifest path.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var q struct {
			Workspace *struct {
				Manifest *struct {
					HostPath string
				} `graphql:"findManifest(format: $format)"`
			} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
		}
		vars := map[string]interface{}{}
		if manifestResolveFlags.Format == "" {
			vars["format"] = (*string)(nil)
		} else {
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
