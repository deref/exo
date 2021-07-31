package main

import (
	"fmt"

	"github.com/deref/exo/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(refreshCmd)
}

var refreshCmd = &cobra.Command{
	Use:   "refresh [refs...]",
	Short: "Refreshes components",
	Long: `Refreshes the state of components.
	
If no components are specified, refreshes all components in the current workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)

		if len(args) == 0 {
			_, err := workspace.RefreshAllComponents(ctx, &api.RefreshAllComponentsInput{})
			return err
		} else {
			// TODO: RefreshComponent should be a bulk operation.
			for _, ref := range args {
				if _, err := workspace.RefreshComponent(ctx, &api.RefreshComponentInput{
					Ref: ref,
				}); err != nil {
					return fmt.Errorf("refreshing %q: %w", ref, err)
				}
			}
			return nil
		}
	},
}
