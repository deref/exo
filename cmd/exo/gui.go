package main

import (
	"os"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(guiCmd)
}

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Opens the exo gui in a web browser",
	Long:  `Opens the exo gui in a web browser.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		browser.Stdout = os.Stderr
		return browser.OpenURL(runState.URL)
	},
}
