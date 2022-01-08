package cli

import (
	"context"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

type controlFunc func(ctx context.Context, ws api.Workspace, refs []string) (jobID string, err error)

func controlComponents(cmd *cobra.Command, args []string, f controlFunc) error {
	ctx := cmd.Context()
	checkOrEnsureServer()
	cl := newClient()
	kernel := cl.Kernel()
	workspace := requireCurrentWorkspace(ctx, cl)

	var jobID string
	var err error
	if len(args) == 0 {
		jobID, err = f(ctx, workspace, nil)
	} else {
		jobID, err = f(ctx, workspace, args)
	}
	if err != nil {
		return err
	}

	return watchJob(ctx, kernel, jobID)
}
