package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().UintVar(timeoutSeconds, "timeout", 0, "The timeout for stopping the process")
}

var restartCmd = &cobra.Command{
	Use:   "restart [ref...]",
	Short: "Restart processes",
	Long: `Restart processes. If not already running, will start them.

If no refs are provided, restarts the entire workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vars := map[string]interface{}{}
		if cmd.Flags().Lookup("timeout").Changed {
			vars["timeoutSeconds"] = timeoutSeconds
		} else {
			vars["timeoutSeconds"] = nil
		}
		return controlComponents(cmd, args, "restartWorkspace", "restartWorkspaceComponents", vars)
	},
}
