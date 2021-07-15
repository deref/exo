package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exitCmd)
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Stop the exo deamon",
	Long:  `Stop the exo deamon process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("TODO: exit command")
	},
}
