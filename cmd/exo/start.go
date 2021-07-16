package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start <ref>",
	Short: "Start a process",
	Long:  `Start a process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		return errors.New("TODO: start command")
	},
}
