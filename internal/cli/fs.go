package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(fsCmd)

	fsCmd.AddCommand(makeHelpSubcmd())
}

var fsCmd = &cobra.Command{
	Hidden: true,
	Use:    "fs",
	Short:  "Filesystem operations.",
	Long:   `Operate on a workspace's file system.`,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}
