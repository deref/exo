package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Hidden: true,
	Use:    "server",
	Short:  "Runs the exo server",
	Long: `Runs the exo server until interrupted.

Prefer the deamonize command for normal operation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		return errors.New("TODO: server command")
	},
}
