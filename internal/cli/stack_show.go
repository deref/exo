package cli

import (
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackShowCmd)
	stackShowCmd.Flags().BoolVar(&stackShowFlags.Recursive, "recursive", false, "include subcomponents")
	stackShowCmd.Flags().BoolVar(&stackShowFlags.Final, "final", false, "fully evaluate to a concrete value")
}

var stackShowFlags struct {
	Recursive bool
	Final     bool
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
				Configuration string `graphql:"configuration(recursive: $recursive, final: $final)"`
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, map[string]any{
			"recursive": stackShowFlags.Recursive,
			"final":     stackShowFlags.Final,
		})
		cmdutil.Show(q.Stack.Configuration)
		return nil
	},
}
