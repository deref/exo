package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	completionCmd.AddCommand(completionUninstallCmd)
}

var completionUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls shell completions",
	Long: `Uninstalls shell completions.

If shell is not specified, removes exo shell completion files from all supported
shells.

After running this command, you must restart your shell. Some shells, such as
zsh, may require additional steps to clear completion caches.`,
	Run: func(cmd *cobra.Command, args []string) {
		shell := ""
		if len(args) > 0 {
			shell = args[0]
		}
		completionUninstall(shell)
	},
}

func completionUninstall(shell string) {
	var shells []string
	if shell == "" {
		shells = completionSupportedShells
	} else {
		shells = []string{shell}
	}

	for _, shell := range shells {
		for _, candidate := range completionPathCandidates(shell) {
			removed := os.Remove(candidate) == nil
			if removed {
				fmt.Println("removed", candidate)
			} else {
				fmt.Println("skipped", candidate)
			}
		}
	}
}

var completionSupportedShells = []string{
	"bash",
	"zsh",
	"fish",
	"powershell",
}
