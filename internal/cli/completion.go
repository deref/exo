package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Hidden:                true,
	Use:                   "completion",
	Short:                 "Generate and install shell completions",
	Long:                  `Generate and install shell completions.`,
	DisableFlagsInUseLine: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		offline = true
		return cmd.Parent().PersistentPreRunE(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
