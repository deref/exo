package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackShowCmd)
}

var stackShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show stack configuration",
	Long:  `Show a stacks's effective configuration.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var q struct {
			Stack *struct {
				Configuration string
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, nil)
		fmt.Println(q.Stack.Configuration)
		return nil
	},
}