package cli

import (
	"os"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
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
		ctx := newContext()
		m, err := loadManifest(ctx, os.Stderr, args[0])
		if err != nil {
			return err
		}
		_, err = hclgen.WriteTo(os.Stdout, m.File)
		return err
	},
}
