package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	workspaceCmd.AddCommand(workspaceDestroyCmd)
}

var workspaceDestroyCmd = &cobra.Command{
	Use:   "destroy [workspace]",
	Short: "Destroys a workspace",
	Long: `Destroys a workspace. If the workspace is not specified, destroys
the workspace for the current working directory.

If the workspace is associated with a current stack, that stack is also
destroyed.`,
	// TODO: Should the files on disk also be destroyed?
	// If yes, recommend `stack destroy` as an alternative and also provide a
	// `workspace forget` command.
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		vars := map[string]any{}
		if len(args) < 1 {
			vars["workspace"] = currentWorkspaceRef()
		} else {
			vars["workspace"] = args[0]
		}
		var m struct {
			Reconciliation struct {
				JobID string
			} `graphql:"destroyWorkspace(workspace: $workspace)"`
		}
		if err := api.Mutate(ctx, svc, &m, vars); err != nil {
			return err
		}
		return watchOwnJob(ctx, m.Reconciliation.JobID)
	},
}
