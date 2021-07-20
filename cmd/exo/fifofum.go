package main

import (
	"github.com/deref/exo/fifofum"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fifofumCmd)
}

var fifofumCmd = &cobra.Command{
	Hidden: true,
	Use:    "fifofum",
	Short:  "Executes and supervises a command with logging to fifos",
	Long:   `Executes and supervises a command with logging to fifos.`,
	Run: func(cmd *cobra.Command, args []string) {
		fifofum.Main("exo fifofum", args)
	},
}
