package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.AddCommand(makeHelpSubcmd())
}

var manifestCmd = &cobra.Command{
	Use:    "manifest",
	Short:  "Manifest tools",
	Long:   `Contains subcommands for working with manifests`,
	Hidden: true,
	Args:   cobra.NoArgs,
}
