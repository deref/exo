package main

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

var initTimeout uint = 0
var timeoutSeconds = &initTimeout

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().UintVar(timeoutSeconds, "timeout", 0, "The timeout for stopping the process")
}

var stopCmd = &cobra.Command{
	Use:   "stop [ref...]",
	Short: "Stop processes",
	Long: `Stop processes.

If no refs are provided, stops the entire workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Lookup("timeout").Changed {
			timeoutSeconds = nil
		}
		return controlComponents(args, func(ctx context.Context, ws api.Workspace, refs []string) (jobID string, err error) {
			if refs == nil {
				output, err := ws.Stop(ctx, &api.StopInput{TimeoutSeconds: timeoutSeconds})
				if output != nil {
					jobID = output.JobID
				}
				return jobID, err
			} else {
				output, err := ws.StopComponents(ctx, &api.StopComponentsInput{
					Refs:           refs,
					TimeoutSeconds: timeoutSeconds,
				})
				if output != nil {
					jobID = output.JobID
				}
				return jobID, err
			}
		})
	},
}
