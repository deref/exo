package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(vaultCmd)

	vaultCmd.AddCommand(makeHelpSubcmd())
}

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Create, inspect, and modify vaults",
	Long:  `Contains subcommands for operating on vaults.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}
