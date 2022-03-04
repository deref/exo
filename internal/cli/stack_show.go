package cli

import (
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackShowCmd)
	stackShowCmd.Flags().BoolVar(&stackShowFlags.Recursive, "recursive", false, "include subcomponents")
}

var stackShowFlags struct {
	Recursive bool
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
				Configuration string `graphql:"configuration(recursive: $recursive)"`
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, map[string]interface{}{
			"recursive": stackShowFlags.Recursive,
		})
		cmdutil.Show(q.Stack.Configuration)
		return nil
	},
}
