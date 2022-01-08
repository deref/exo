package cli

import (
	"os"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestExpandCmd)
}

var manifestExpandCmd = &cobra.Command{
	Use:   "expand <manifest>",
	Short: "Loads a manifest file and prints it in expanded form.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		m, err := loadManifest(ctx, os.Stderr, args[0])
		if err != nil {
			return err
		}

		expanded := exohcl.RewriteManifest(&exohcl.Expand{Context: ctx}, m)
		hclgen.WriteTo(os.Stdout, expanded)
		return nil
	},
}
