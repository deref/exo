package main

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
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
		kernel := cl.Kernel()
		workspace := requireWorkspace(ctx, cl)

		var input api.RefreshComponentsInput
		if len(args) > 0 {
			input.Refs = args
		}

		output, err := workspace.RefreshComponents(ctx, &input)
		if err != nil {
			return fmt.Errorf("refreshing: %w", err)
		}

		return watchJob(ctx, kernel, output.JobID)
	},
}
