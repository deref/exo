package main

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd)
}

var restartCmd = &cobra.Command{
	Use:   "restart [ref...]",
	Short: "Restart processes",
	Long: `Restart processes. If not already running, will start them.

If no refs are provided, restarts the entire workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return controlComponents(args, func(ctx context.Context, ws api.Workspace, refs []string) (jobID string, err error) {
			if refs == nil {
				output, err := ws.Restart(ctx, &api.RestartInput{})
				if output != nil {
					jobID = output.JobID
				}
				return jobID, err
			} else {
				output, err := ws.RestartComponents(ctx, &api.RestartComponentsInput{
					Refs: refs,
				})
				if output != nil {
					jobID = output.JobID
				}
				return jobID, err
			}
		})
	},
}
