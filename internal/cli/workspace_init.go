package cli

import (
	"fmt"

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
				ID string
			} `graphql:"newWorkspace(root: $root)"`
		}
		if err := client.Mutate(ctx, &m, map[string]interface{}{
			"root": root,
		}); err != nil {
			return err
		}
		fmt.Println(m.Workspace.ID)
		return nil
	},
}
