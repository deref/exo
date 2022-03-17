package cli

import (
	"github.com/spf13/cobra"
)

var initTimeout uint = 0
var timeoutSeconds = &initTimeout

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().UintVar(timeoutSeconds, "timeout", 0, "The timeout for stopping the process")
}

var stopCmd = &cobra.Command{
	Use:   "stop [ref...]",
	Short: "Stop processes",
	Long: `Stop processes.

If no refs are provided, stops the entire workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vars := map[string]any{}
		if cmd.Flags().Lookup("timeout").Changed {
			vars["timeoutSeconds"] = timeoutSeconds
		} else {
			vars["timeoutSeconds"] = nil
		}
		return controlComponents(cmd, args, "stopWorkspace", "stopWorkspaceComponents", vars)
	},
}
