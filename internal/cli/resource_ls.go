package cli

import (
	"fmt"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceLSCmd)
	resourceLSCmd.Flags().BoolVarP(&resourceLSFlags.All, "all", "a", false, "Alias for --scope=all")
	resourceLSCmd.Flags().StringVar(&resourceLSFlags.Scope, "scope", "", "component, stack, project, or all")
	resourceLSCmd.Flags().StringVar(&resourceLSFlags.Component, "component", "", "")
}

var resourceLSFlags struct {
	All       bool
	Scope     string
	Component string
}

var resourceLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "List resources",
	Long:  `Lists resources in the given scope.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		scope := resourceLSFlags.Scope
		if resourceLSFlags.All {
			if scope == "" {
				scope = "all"
			} else if scope != "all" {
				return fmt.Errorf("--all conflicts with --scope=%q", scope)
			}
		}
		if resourceLSFlags.Component != "" {
			if scope == "" {
				scope = "component"
			} else if scope != "component" {
				return fmt.Errorf("--component conflicts with --scope=%q", scope)
			}
		}
		if scope == "" {
			scope = "stack"
		}

		type resourceFragment struct {
			IRI     string
			Project *struct {
				DisplayName string
			}
			Stack *struct {
				Name string
			}
			Component *struct {
				Name string
			}
		}
		var resources []resourceFragment
		var columns []string

		switch scope {
		case "all":
			var q struct {
				Resources []resourceFragment `graphql:"allResources"`
			}
			if err := client.Query(ctx, &q, nil); err != nil {
				return err
			}
			resources = q.Resources
			columns = []string{"IRI", "PROJECT", "STACK", "COMPONENT"}

		case "project":
			var q struct {
				Workspace *struct {
					Project struct {
						Resources []resourceFragment
					}
				} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
			}
			mustQueryWorkspace(ctx, client, &q, nil)
			resources = q.Workspace.Project.Resources
			columns = []string{"IRI", "STACK", "COMPONENT"}

		case "stack":
			var q struct {
				Stack *struct {
					ID        string
					Resources []resourceFragment
				} `graphql:"stackByRef(ref: $currentStack)"`
			}
			mustQueryStack(ctx, client, &q, nil)
			resources = q.Stack.Resources
			columns = []string{"IRI", "COMPONENT"}

		case "component":
			var q struct {
				Stack *struct {
					ID        string
					Resources []resourceFragment
				} `graphql:"stackByRef(ref: $currentStack)"`
			}
			mustQueryStack(ctx, client, &q, nil)
			resources = q.Stack.Resources
			columns = []string{"IRI"}

		default:
			return fmt.Errorf("unknown scope: %q", resourceLSFlags.Scope)
		}

		w := cmdutil.NewTableWriter(columns...)
		for _, resource := range resources {
			data := map[string]string{
				"IRI": resource.IRI,
			}
			if resource.Project != nil {
				data["PROJECT"] = resource.Project.DisplayName
			}
			if resource.Stack != nil {
				data["STACK"] = resource.Stack.Name
			}
			if resource.Component != nil {
				data["COMPONENT"] = resource.Component.Name
			}
			values := make([]string, len(columns))
			for i, column := range columns {
				values[i] = data[column]
			}
			w.WriteRow(values...)
		}
		w.Flush()
		return nil
	},
}
