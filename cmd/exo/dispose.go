package main

import (
	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(disposeCmd)
}

var disposeCmd = &cobra.Command{
	Use:    "dispose [ref ...]",
	Short:  "Disposes components",
	Long:   "Disposes components.",
	Hidden: true, // This command is only really useful for testing controllers.
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		_, err := workspace.DisposeComponents(ctx, &api.DisposeComponentsInput{
			Refs: args,
		})
		return err
	},
}
