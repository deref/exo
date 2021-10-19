package cli

import (
	"context"

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
		return controlComponents(args, func(ctx context.Context, ws api.Workspace, refs []string) (jobID string, err error) {
			var input api.RefreshComponentsInput
			if len(args) > 0 {
				input.Refs = args
			}
			output, err := ws.RefreshComponents(ctx, &input)
			if err != nil {
				return "", err
			}
			return output.JobID, nil
		})
	},
}
