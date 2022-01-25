package cli

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new <type> <name> [args]",
	Short: "Creates a new component",
	Long: `Creates a new component of a given type with a given name.  Each
component type may define its own syntax for flags and positional arguments.

To learn about specific types - for example Processes - consult each type's
help page:

exo help new process
`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func createComponent(ctx context.Context, name, typ string, spec interface{}) error {
	var m struct {
		Reconciliation struct {
			Component struct {
				ID string
			}
			Job struct {
				ID string
			}
		} `graphql:"createComponent(stack: $stack, name: $name, type: $type, spec: $spec)"`
	}
	if err := api.Mutate(ctx, svc, &m, map[string]interface{}{
		"stack": currentStackRef(),
		"name":  name,
		"type":  typ,
		"spec":  jsonutil.MustMarshalString(spec),
	}); err != nil {
		return err
	}
	fmt.Println("Component-ID:", m.Reconciliation.Component.ID)
	return watchJob(ctx, m.Reconciliation.Job.ID)
}
