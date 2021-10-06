package main

import (
	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:   "rename <ref> <new-name>",
	Short: "Rename component",
	Long:  "Rename components.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		_, err := workspace.RenameComponent(ctx, &api.RenameComponentInput{
			Ref:  args[0],
			Name: args[1],
		})
		return err
	},
}
