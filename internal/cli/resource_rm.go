package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceRmCmd)
}

var resourceRmCmd = &cobra.Command{
	Use:   "rm <ref>", // TODO: Variadic.
	Short: "Remove a resource",
	Long:  "Dispose a resource via its provider, then stop tracking it.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return sendMutation(ctx, "disposeResource", map[string]interface{}{
			"ref": args[0],
		})
	},
}
