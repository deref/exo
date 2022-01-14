package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceAdoptCmd)
	resourceAdoptCmd.Flags().StringVar(&resourceAdoptFlags.Owner, "owner", "", "component, stack, project, or none")
	resourceAdoptCmd.Flags().StringVar(&resourceAdoptFlags.Component, "component", "", "Set owner to component")
}

var resourceAdoptFlags struct {
	Owner     string
	Component string
}

var resourceAdoptCmd = &cobra.Command{
	Use:   "adopt <iri>",
	Short: "Take ownership of a resource",
	Long: `Take ownership of a resource.

If --owner is specified, the resource's new owner will be set to the current
entity of that type relative to the current workspace.

--component implies --owner=component

If --owner is not specified, the new owner will be the first available of
current stack or current project.

--owner=none will track orphaned resources.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		iri := args[0]

		ownerType := resourceAdoptFlags.Owner
		if resourceAdoptFlags.Component != "" {
			if ownerType == "" {
				ownerType = "component"
			} else if ownerType != "component" {
				return fmt.Errorf("--component conflicts with --owner=%q", ownerType)
			}
		}

		var m struct {
			Resource struct {
				IRI string
			} `graphql:"adoptResource(iri: $iri, workspace: $workspace, ownerType: $ownerType, component: $component)"`
		}
		vars := map[string]interface{}{
			"iri": iri,
		}
		switch ownerType {
		case "component":
			vars["ownerType"] = "Component"
			vars["workspace"] = currentWorkspaceRef()
			vars["component"] = resourceAdoptFlags.Component
		case "stack":
			vars["ownerType"] = "Stack"
			vars["workspace"] = currentWorkspaceRef()
			vars["component"] = (*string)(nil)
		case "project":
			vars["ownerType"] = "Project"
			vars["workspace"] = currentWorkspaceRef()
			vars["component"] = (*string)(nil)
		case "none":
			vars["ownerType"] = (*string)(nil)
			vars["workspace"] = (*string)(nil)
			vars["component"] = (*string)(nil)
		case "":
			vars["ownerType"] = (*string)(nil)
			vars["workspace"] = currentWorkspaceRef()
			vars["component"] = (*string)(nil)
		default:
			return fmt.Errorf("unexpected value for --owner: %q", ownerType)
		}
		return cl.Mutate(ctx, &m, vars)
	},
}
