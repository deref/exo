package cli

import (
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/spf13/cobra"
)

func controlComponents(cmd *cobra.Command, args []string, workspaceMutation string, componentsMutation string, vars map[string]interface{}) error {
	ctx := cmd.Context()
	checkOrEnsureServer()
	kernel := newClient().Kernel()

	cl, shutdown := dialGraphQL(ctx)
	defer shutdown()

	// TODO: It would be nice to have generated mutation methods.
	var mutation string
	vars = jsonutil.Merge(map[string]interface{}{
		"workspace": currentWorkspaceRef(),
	}, vars)
	if len(args) == 0 {
		mutation = workspaceMutation
	} else {
		mutation = componentsMutation
		vars["components"] = args
	}

	jobID, err := cl.StartTask(ctx, mutation, vars)
	if err != nil {
		return err
	}

	return watchJob(ctx, kernel, jobID)
}
