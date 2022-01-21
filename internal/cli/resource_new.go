package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceNewCmd)
	resourceNewCmd.Flags().BoolVar(&resourceNewFlags.Adopt, "adopt", false, "do not create external resource")
	resourceNewCmd.Flags().StringVar(&resourceNewFlags.Owner, "owner", "", "component, stack, project, or none")
	resourceNewCmd.Flags().StringVar(&resourceNewFlags.Component, "component", "", "set owner to component")
}

var resourceNewFlags struct {
	Adopt     bool
	Owner     string
	Component string
}

var resourceNewCmd = &cobra.Command{
	Use:   "new <type> [property...]",
	Short: "Create a new resource",
	Long: `Create a new resource.

Desired resource state is specified by a model string with a format defined by
the resource type. Many resource types use JSON for models, which can be
conveniently specified as property arguments on the command line.  See 'exo
json --help' for property syntax. If no properties are provided, the model
string is read from stdin, or entered interactively.

If --adopt is specified, a new resource record will be created to track an
existing external resource, but the resource controller will be instructed not
to create a new external resource. Instead, the controller will only read a
refreshed resource model. The resource model given to the 'resource new'
command invocation must include enough information for the controller to
uniquely identify the existing external resource.

An owner will be assigned to the new resource, defaulting to the first
available of the current stack or the current project. If --owner is provided,
the resource's owner will be set to the current entity of that type relative to
the current workspace.

To set the owner to a component, supply '--component=ref', which implies
'--owner=component'.

'--owner=none' allows tracking of orphaned resources.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		typ := args[0]
		props := args[1:]

		var model string
		var err error
		if len(props) > 0 {
			var obj map[string]interface{}
			obj, err = cmdutil.ArgsToJsonObject(props)
			if err == nil {
				model = jsonutil.MustMarshalString(obj)
			}
		} else if term.IsInteractive() {
			model, err = term.EditString("resource.*.txt", "")
		} else {
			var bs []byte
			bs, err = ioutil.ReadAll(os.Stdin)
			model = string(bs)
		}
		if err != nil {
			return err
		}
		if model == "" {
			return errors.New("empty resource model")
		}

		ownerType := resourceNewFlags.Owner
		if resourceNewFlags.Component != "" {
			if ownerType == "" {
				ownerType = "component"
			} else if ownerType != "component" {
				return fmt.Errorf("--component conflicts with --owner=%q", ownerType)
			}
		}

		var m struct {
			Resource struct {
				ID     string
				TaskID *string
			} `graphql:"newResource(type: $type, model: $model, workspace: $workspace, ownerType: $ownerType, component: $component, adopt: $adopt)"`
		}
		vars := map[string]interface{}{
			"type":  typ,
			"model": model,
			"adopt": resourceNewFlags.Adopt,
		}
		switch ownerType {
		case "component":
			vars["ownerType"] = "Component"
			vars["workspace"] = currentWorkspaceRef()
			vars["component"] = resourceNewFlags.Component
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
		if err := api.Mutate(ctx, svc, &m, vars); err != nil {
			return err
		}
		fmt.Println("resource:", m.Resource.ID)
		if rootPersistentFlags.Async {
			return nil
		}
		return watchOwnJob(ctx, *m.Resource.TaskID)
	},
}
