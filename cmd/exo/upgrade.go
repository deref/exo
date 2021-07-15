package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade exo",
	Long:  `Upgrade exo to the latest version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("TODO: upgrade command")
	},
}
