package cli

import (
	"github.com/deref/exo/internal/exod"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Hidden: true,
	Use:    "server",
	Short:  "Runs the exo server",
	Long: `Runs the exo server until interrupted.

Prefer the daemonize command for normal operation.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		svr := &exod.Server{
			RedirectCrashDumps: !logToStderr(),
		}
		svr.Run(ctx)
	},
}
