package main

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

var stopNow = false

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVar(&stopNow, "stop-now", false, "Attempts to stop the process immediately")
}

var stopCmd = &cobra.Command{
	Use:   "stop [ref...]",
	Short: "Stop processes",
	Long: `Stop processes.

If no refs are provided, stops the entire workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return controlComponents(args, func(ctx context.Context, ws api.Workspace, refs []string) (jobID string, err error) {
			if refs == nil {
				output, err := ws.Stop(ctx, &api.StopInput{StopNow: stopNow})
				if output != nil {
					jobID = output.JobID
				}
				return jobID, err
			} else {
				output, err := ws.StopComponents(ctx, &api.StopComponentsInput{
					Refs:    refs,
					StopNow: stopNow,
				})
				if output != nil {
					jobID = output.JobID
				}
				return jobID, err
			}
		})
	},
}
