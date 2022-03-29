package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackDestroyCmd)
}

var stackDestroyCmd = &cobra.Command{
	Use:   "destroy [stack]",
	Short: "Destroys a stack",
	Long: `Destroys a stack. If the stack is not specified, destroys
the current stack.

Destroying a stack also destroys all components in that stack.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		vars := map[string]any{}
		if len(args) < 1 {
			vars["stack"] = currentStackRef()
		} else {
			vars["stack"] = args[0]
		}
		var m struct {
			Reconciliation struct {
				JobID string
			} `graphql:"destroyStack(ref: $stack)"`
		}
		if err := api.Mutate(ctx, svc, &m, vars); err != nil {
			return err
		}
		return watchOwnJob(ctx, m.Reconciliation.JobID)
	},
}
