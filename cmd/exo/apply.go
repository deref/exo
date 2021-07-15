package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVar(&applyFlags.Format, "format", "", "exo, procfile, or compose")
}

var applyFlags struct {
	Format string
}

var applyCmd = &cobra.Command{
	Use:   "apply [flags] <config-file>",
	Short: "Applies a config to the current project",
	Long:  `Applies a config to the current project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("TODO: apply command")
	},
}
