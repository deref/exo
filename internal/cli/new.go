package cli

import (
	"context"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
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

func createComponent(ctx context.Context, name, typ string, spec any) error {
	var m struct {
		Reconciliation struct {
			Component struct {
				ID string `json:"id"`
			}
			JobID string
		} `graphql:"createComponent(stack: $stack, name: $name, type: $type, spec: $spec)"`
	}
	if err := api.Mutate(ctx, svc, &m, map[string]any{
		"stack": currentStackRef(),
		"name":  name,
		"type":  typ,
		"spec":  jsonutil.MustMarshalString(spec),
	}); err != nil {
		return err
	}
	cmdutil.PrintCueStruct(m.Reconciliation.Component)
	return watchOwnJob(ctx, m.Reconciliation.JobID)
}
