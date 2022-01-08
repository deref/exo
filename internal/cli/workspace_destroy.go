package cli

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	workspaceCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy [workspace]",
	Short: "Deletes a workspace",
	Long: `Deletes a workspace. If the workspace is not specified, deletes
the workspace for the current working directory.

Deleting a workspace also deletes all resources in that workspace.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()
		cl := newClient()
		kernel := cl.Kernel()
		var workspace api.Workspace
		if len(args) < 1 {
			workspace = requireCurrentWorkspace(ctx, cl)
		} else {
			ref := args[0]
			var err error
			workspace, err = resolveWorkspace(ctx, cl, ref)
			if err != nil {
				return fmt.Errorf("resolving workspace: %w", err)
			}
			if workspace == nil {
				return fmt.Errorf("unresolved workspace ref: %q", ref)
			}
		}
		output, err := workspace.Destroy(ctx, &api.DestroyInput{})
		if err != nil {
			return err
		}
		return watchJob(ctx, kernel, output.JobID)
	},
}
