package main

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	workspaceCmd.AddCommand(workspaceInitCmd)
}

var workspaceInitCmd = &cobra.Command{
	Use:   "init [root]",
	Short: "Creates a workspace",
	Long: `Creates a workspace. If root is not provided, the new workspace
will be rooted at the current working directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		var root string
		if len(args) < 1 {
			root = cmdutil.MustGetwd()
		} else {
			root = args[0]
		}
		output, err := cl.Kernel().CreateWorkspace(ctx, &api.CreateWorkspaceInput{
			Root: root,
		})
		if err != nil {
			cmdutil.Fatalf("creating workspace: %w", err)
		}
		fmt.Println(output.ID)
		return nil
	},
}
