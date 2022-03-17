package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackNewCmd)
	stackNewCmd.Flags().StringVar(&stackNewFlags.Name, "name", "", "Name of stack")
}

var stackNewFlags struct {
	Name    string
	Cluster string
	Detatch bool
}

var stackNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new stack",
	Long: `Create a new stack.

Associates the new stack with the current workspace.

If a name is not provided, the stack's name will be set to its generated id.

If --cluster is not specified, the stack will be created in the default
cluster.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		vars := map[string]any{
			"workspace": currentWorkspaceRef(),
		}
		if cmd.Flags().Lookup("name").Changed {
			vars["name"] = stackNewFlags.Name
		} else {
			vars["name"] = (*string)(nil)
		}
		if cmd.Flags().Lookup("cluster").Changed {
			vars["cluster"] = stackNewFlags.Cluster
		} else {
			vars["cluster"] = (*string)(nil)
		}
		var m struct {
			Stack struct {
				ID string `json:"id"`
			} `graphql:"createStack(name: $name, workspace: $workspace, cluster: $cluster)"`
		}
		if err := api.Mutate(ctx, svc, &m, vars); err != nil {
			return err
		}
		cmdutil.PrintCueStruct(m.Stack)
		return nil
	},
}
