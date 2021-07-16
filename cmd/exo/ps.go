package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lists defined processes.",
	Long:  `Describes defined processes and their statuses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		return errors.New("TODO: ps command")
	},
}
