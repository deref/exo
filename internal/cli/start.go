package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start [ref...]",
	Short: "Start processes",
	Long: `Start processes.

If no refs are provided, starts the entire workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return controlComponents(cmd, args, "startWorkspace", "startWorkspaceComponents", nil)
	},
}
