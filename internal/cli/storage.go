package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(storageCmd)

	storageCmd.AddCommand(makeHelpSubcmd())
}

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Create, inspect, and modify stores",
	Long:  `Contains subcommands for operating on stores.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}
