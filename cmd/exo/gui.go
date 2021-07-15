package main

import (
	"errors"

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
		return errors.New("TODO: open command")
	},
}
