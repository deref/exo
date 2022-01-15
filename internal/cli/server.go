package cli

import (
	"github.com/deref/exo/internal/exod"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Hidden:             true,
	Use:                "server",
	Short:              "Runs the exo server",
	DisableFlagParsing: true,
	Long: `Runs the exo server until interrupted.

Prefer the daemonize command for normal operation.`,
	Args: cobra.ArbitraryArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		offline = true
		return cmd.Parent().PersistentPreRunE(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		exod.Main(ctx)
	},
}
