package main

import (
	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm [ref ...]",
	Short: "Remove components",
	Long:  "Remove components.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)

		_, err := workspace.DeleteComponents(ctx, &api.DeleteComponentsInput{
			Refs: args,
		})
		return err
	},
}
