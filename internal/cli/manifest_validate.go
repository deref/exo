package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestValidateCmd)
}

var manifestValidateCmd = &cobra.Command{
	Use:   "validate <manifest>",
	Short: "Loads and validates a manifest file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		_, err := loadManifest(ctx, os.Stdout, args[0])
		return err
	},
}
