package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resourceCmd)

	resourceCmd.AddCommand(makeHelpSubcmd())
}

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Create, inspect, and modify resources",
	Long:  `Contains subcommands for operating on resources.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}
