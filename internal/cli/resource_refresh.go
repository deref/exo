package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceRefreshCmd)
}

var resourceRefreshCmd = &cobra.Command{
	Use:   "refresh <ref>", // TODO: Variadic.
	Short: "Refresh a resource",
	Long:  "Refresh a resource's modelby reading external state.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return sendMutation(ctx, "refreshResource", map[string]any{
			"ref": args[0],
		})
	},
}
