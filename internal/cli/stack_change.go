package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackChangeCmd)
}

var stackChangeCmd = &cobra.Command{
	Use:   "change <ref>",
	Short: "Change current stack",
	Long: `Change the current workspace's current stack.

Use the ref '-' to clear the current stack.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var stack *string
		if args[0] != "-" {
			stack = &args[0]
		}

		var m struct {
			Stack struct {
				ID string
			} `graphql:"setWorkspaceStack(workspace: $workspace, stack: $stack)"`
		}
		return api.Mutate(ctx, svc, &m, map[string]interface{}{
			"workspace": currentWorkspaceRef(),
			"stack":     stack,
		})
	},
}
