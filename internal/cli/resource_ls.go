package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceLSCmd)
	resourceLSCmd.Flags().BoolVarP(&resourceLSFlags.All, "all", "a", false, "Alias for --scope=all")
	resourceLSCmd.Flags().StringVar(&resourceLSFlags.Scope, "scope", "stack", "stack, project, all")
}

var resourceLSFlags struct {
	All   bool
	Scope string
}

var resourceLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "List resources",
	Long:  `Lists resources in the given scope.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		if resourceLSFlags.All {
			if cmd.Flags().Lookup("scope").Changed {
				return errors.New("--all and --scope are mutually exclusive")
			}
			resourceLSFlags.Scope = "all"
		}

		type resourceFragment struct {
			IRI string
		}
		var resources []resourceFragment

		switch resourceLSFlags.Scope {
		case "stack":
			var q struct {
				Stack *struct {
					ID        string
					Resources []resourceFragment
				} `graphql:"stackByRef(ref: $currentStack)"`
			}
			mustQueryStack(ctx, cl, &q, nil)
			resources = q.Stack.Resources

		case "all":
			var q struct {
				Resources []resourceFragment `graphql:"allResources"`
			}
			if err := cl.Query(ctx, &q, nil); err != nil {
				return err
			}
			resources = q.Resources

		case "project":
			var q struct {
				Workspace *struct {
					Project struct {
						Resources []resourceFragment
					}
				} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
			}
			mustQueryWorkspace(ctx, cl, &q, nil)
			resources = q.Workspace.Project.Resources

		default:
			return fmt.Errorf("unknown scope: %q", resourceLSFlags.Scope)
		}

		for _, resource := range resources {
			fmt.Println(resource.IRI)
		}
		return nil
	},
}
