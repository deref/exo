package cli

import (
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exitCmd)
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Stop the exo daemon",
	Long:  `Stop the exo daemon process.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.NoDaemon {
			cmdutil.Fatalf("daemon disabled by config")
		}
		ctx := cmd.Context()
		// TODO: Fail gracefully if the daemon is already stopped, or if
		// it exists before the response comes back.
		var resp struct{}
		return svc.Do(ctx, `mutation { stopDaemon }`, nil, &resp)
	},
}
