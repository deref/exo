package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm <ref>",
	Short: "Remove component",
	// TODO: Update manifest?
	Long: "Remove component from the current stack.",
	Args: cobra.ExactArgs(1), // TODO: Variadic.
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ref := args[0]
		var m struct {
			Reconciliation struct {
				Job struct {
					ID string
				}
			} `graphql:"disposeComponent(stack: $stack, ref: $component)"`
		}
		if err := api.Mutate(ctx, svc, &m, map[string]interface{}{
			"stack":     currentStackRef(),
			"component": ref,
		}); err != nil {
			return err
		}
		return watchJob(ctx, m.Reconciliation.Job.ID)
	},
}
