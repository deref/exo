package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestFormatCmd)
}

var manifestFormatCmd = &cobra.Command{
	Use:   "format <manifest>",
	Short: "Loads a manifest file and reformats it to standard out.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// XXX
		return nil
	},
}
