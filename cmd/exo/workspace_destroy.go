package main

import (
	"github.com/deref/exo/core/api"
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
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		var workspace api.Workspace
		if len(args) < 1 {
			workspace = requireWorkspace(ctx, cl)
		} else {
			workspace = cl.GetWorkspace(args[0])
		}
		_, err := workspace.Destroy(ctx, &api.DestroyInput{})
		return err
	},
}
