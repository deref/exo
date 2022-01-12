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
		checkOrEnsureServer() // XXX should not ensure the server if we're trying to exit it!

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var resp struct{}
		return cl.Run(ctx, `mutation { stopDaemon }`, &resp, nil)
	},
}
