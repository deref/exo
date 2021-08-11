package main

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:    "build [ref ...]",
	Short:  "Build components",
	Long:   "Build components.",
	Hidden: true, // Not sure if we want to expose build functionality yet.
	RunE: func(cmd *cobra.Command, args []string) error {
		return controlComponents(args, func(ctx context.Context, ws api.Workspace, refs []string) (jobID string, err error) {
			if refs == nil {
				output, err := ws.Build(ctx, &api.BuildInput{})
				if output != nil {
					jobID = output.JobID
				}
				return jobID, err
			} else {
				output, err := ws.BuildComponents(ctx, &api.BuildComponentsInput{
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
