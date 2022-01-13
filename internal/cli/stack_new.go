package cli

import (
	"fmt"

	"github.com/shurcooL/graphql"
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
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		vars := map[string]interface{}{
			"workspace": graphql.String(currentWorkspaceRef()),
		}
		if cmd.Flags().Lookup("name").Changed {
			vars["name"] = graphql.String(stackNewFlags.Name)
		} else {
			vars["name"] = (*graphql.String)(nil)
		}
		if cmd.Flags().Lookup("cluster").Changed {
			vars["cluster"] = graphql.String(stackNewFlags.Cluster)
		} else {
			vars["cluster"] = (*graphql.String)(nil)
		}
		var m struct {
			Stack struct {
				ID string
			} `graphql:"newStack(name: $name, workspace: $workspace, cluster: $cluster)"`
		}
		if err := cl.Mutate(ctx, &m, vars); err != nil {
			return err
		}
		fmt.Println(m.Stack.ID)
		return nil
	},
}
