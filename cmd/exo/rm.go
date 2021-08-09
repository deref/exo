package main

import (
	"fmt"

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

		// TODO: Bulk delete operation.
		for _, ref := range args {
			_, err := workspace.DeleteComponent(ctx, &api.DeleteComponentInput{
				Ref: ref,
			})
			if err != nil {
				return fmt.Errorf("deleting %q: %w", ref, err)
			}
		}
		return nil
	},
}
