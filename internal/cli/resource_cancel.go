package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceCancelCmd)
}

var resourceCancelCmd = &cobra.Command{
	Use:   "cancel <ref>", // TODO: Variadic.
	Short: "Cancel a resource operation",
	Long: `Cancel an operation being performed by a resource controller and frees
the lock on the resource.

Makes no attempt to cleanup or undo external or asynchronous effects initiated
by a resource controller.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ref := args[0]
		var resp struct{}
		return svc.Do(ctx, `
			mutation ($ref: String!) {
				cancelResourceOperation(ref: $ref) {
					__typename
				}
			}
		`, map[string]interface{}{
			"ref": ref,
		}, &resp)
	},
}
