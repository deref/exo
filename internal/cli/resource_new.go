package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceNewCmd)
	resourceNewCmd.Flags().BoolVar(&resourceNewFlags.Adopt, "adopt", false, "do not create external resource")
	resourceNewCmd.Flags().StringVar(&resourceNewFlags.Owner, "owner", "stack", "component, stack, project, or none")
	resourceNewCmd.Flags().StringVar(&resourceNewFlags.Component, "component", "", "component ref")
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

Desired resource state is specified by a JSON model.  See 'exo json --help' for
property syntax. If no properties are provided, the model string is read from
stdin, or entered interactively.

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

		var model map[string]any
		var err error
		if len(props) > 0 {
			model, err = cmdutil.ArgsToJsonObject(props)
			if err != nil {
				return err
			}
		} else {
			var modelString string
			if isInteractive() {
				modelString, err = term.EditString("resource.*.json", "")
			} else {
				var bs []byte
				bs, err = ioutil.ReadAll(os.Stdin)
				modelString = string(bs)
			}
			if err != nil {
				return err
			}
			modelString = strings.TrimSpace(modelString)
			if modelString == "" {
				return errors.New("empty resource model")
			}
			err = jsonutil.UnmarshalString(modelString, &model)
			if err != nil {
				return err
			}
		}

		ownerType := resourceNewFlags.Owner
		if resourceNewFlags.Component != "" {
			if !cmd.Flag("owner").Changed {
				ownerType = "component"
			} else if ownerType != "component" {
				return fmt.Errorf("--component conflicts with --owner=%q", ownerType)
			}
		}

		var m struct {
			Resource struct {
				ID     string
				TaskID *string
			} `graphql:"createResource(type: $type, model: $model, project: $project, stack: $stack, component: $component, adopt: $adopt)"`
		}
		vars := map[string]any{
			"type":  typ,
			"model": scalars.JSONObject(model),
			"adopt": resourceNewFlags.Adopt,
		}
		switch ownerType {
		case "component":
			vars["project"] = currentProjectRef()
			vars["stack"] = currentStackRef()
			vars["component"] = resourceNewFlags.Component
		case "stack":
			vars["project"] = currentProjectRef()
			vars["stack"] = currentStackRef()
			vars["component"] = (*string)(nil)
		case "project":
			vars["project"] = currentProjectRef()
			vars["stack"] = (*string)(nil)
			vars["component"] = (*string)(nil)
		case "none":
			vars["project"] = (*string)(nil)
			vars["stack"] = (*string)(nil)
			vars["component"] = (*string)(nil)
		default:
			return fmt.Errorf("unexpected value for --owner: %q", ownerType)
		}
		if err := api.Mutate(ctx, svc, &m, vars); err != nil {
			return err
		}
		cmdutil.PrintCueStruct(map[string]any{
			"id": m.Resource.ID,
		})
		if rootPersistentFlags.Async {
			return nil
		}
		return watchOwnJob(ctx, *m.Resource.TaskID)
	},
}
