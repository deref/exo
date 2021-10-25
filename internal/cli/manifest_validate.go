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
		_, err := loadManifest(args[0])
		return writeManifestError(os.Stdout, err)
	},
}
