package cli

import (
	"github.com/deref/exo/internal/supervise"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(superviseCmd)
}

var superviseCmd = &cobra.Command{
	Hidden: true,
	Use:    "supervise",
	Short:  "Supervises a command",
	Long: `Executes a command, supervises its execution, and redirects stdout/stderr to syslog.

This is an internal use command. See the supervise package implementation for usage details.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		offline = true
		return cmd.Parent().PersistentPreRunE(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		supervise.Main()
	},
}
