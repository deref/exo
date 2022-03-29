package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm <refs...>",
	Short: "Remove components",
	// TODO: Update manifest?
	Long: "Remove components from the current stack.",
	Args: cobra.ExactArgs(1), // TODO: Variadic.
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		refs := args
		var m struct {
			Reconciliation struct {
				JobID string
			} `graphql:"destroyComponents(stack: $stack, refs: $components)"`
		}
		if err := api.Mutate(ctx, svc, &m, map[string]any{
			"stack":      currentStackRef(),
			"components": refs,
		}); err != nil {
			return err
		}
		return watchJob(ctx, m.Reconciliation.JobID)
	},
}
