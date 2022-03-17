package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	workspaceCmd.AddCommand(workspaceInitCmd)
}

var workspaceInitCmd = &cobra.Command{
	Use:   "init [root]",
	Short: "Creates a workspace",
	Long: `Creates an empty workspace.

Outputs the ID of the newly created workspace.

If root is not provided, the new workspace will be rooted at the current
working directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var root string
		if len(args) < 1 {
			root = cmdutil.MustGetwd()
		} else {
			root = args[0]
		}
		var m struct {
			Workspace struct {
				ID string `json:"id"`
			} `graphql:"createWorkspace(root: $root)"`
		}
		if err := api.Mutate(ctx, svc, &m, map[string]any{
			"root": root,
		}); err != nil {
			return err
		}
		cmdutil.PrintCueStruct(m.Workspace)
		return nil
	},
}
