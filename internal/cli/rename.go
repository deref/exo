package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:   "rename <ref> <new-name>",
	Short: "Rename component",
	Long:  "Rename components.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ref := args[0]
		newName := args[1]
		return sendMutation(ctx, "updateComponent", map[string]any{
			"stack":   currentStackRef(),
			"ref":     ref,
			"newName": newName,
		})
	},
}
