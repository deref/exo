package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(refreshCmd)
}

var refreshCmd = &cobra.Command{
	Use:   "refresh [refs...]",
	Short: "Refreshes components",
	Long: `Refreshes the state of components.
	
If no components are specified, refreshes all components in the current workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return controlComponents(cmd, args, "refreshWorkspace", "refreshWorkspaceComponents", nil)
	},
}
