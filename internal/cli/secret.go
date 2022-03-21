package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.AddCommand(makeHelpSubcmd())
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Create, inspect, and modify secrets",
	Long:  `Contains subcommands for operating on secrets.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}
