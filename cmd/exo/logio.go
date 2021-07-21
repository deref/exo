package main

import (
	"github.com/deref/exo/logio"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logioCmd)
}

var logioCmd = &cobra.Command{
	Hidden: true,
	Use:    "logio",
	Short:  "Supervises a command with logging to syslog",
	Long: `Executes a command, supervises its execution, and redirects stdout/stderr to syslog.

This is an internal use command. See the logio package implementation for usage details.`,
	Run: func(cmd *cobra.Command, args []string) {
		logio.Main("exo logio", args)
	},
}
