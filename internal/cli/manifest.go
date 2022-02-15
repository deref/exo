package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.AddCommand(makeHelpSubcmd())
}

func addManifestFormatFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "format", "", "exo, exohcl, compose, procfile")
}

var manifestCmd = &cobra.Command{
	Use:    "manifest",
	Short:  "Manifest tools",
	Long:   `Contains subcommands for working with manifests`,
	Hidden: true,
	Args:   cobra.NoArgs,
}
