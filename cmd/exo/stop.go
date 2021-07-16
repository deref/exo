package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop <ref>",
	Short: "Stop a process",
	Long:  `Stop a process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		return errors.New("TODO: stop command")
	},
}
