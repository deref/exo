package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:    "build [ref ...]",
	Short:  "Build components",
	Long:   "Build components.",
	Hidden: true, // Not sure if we want to expose build functionality yet.
	RunE: func(cmd *cobra.Command, args []string) error {
		return controlComponents(cmd, args, "buildWorkspace", "buildWorkspaceComponents", nil)
	},
}
