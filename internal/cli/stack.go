package cli

import (
	"context"
	"reflect"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stackCmd)

	stackCmd.AddCommand(makeHelpSubcmd())
}

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Create, inspect and modify stacks",
	Long: `Contains subcommands for operating on stacks.

If no subcommand is given, describes the current stack of the current
workspace.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var q struct {
			Stack *struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, nil)

		cmdutil.PrintCueStruct(q.Stack)
		return nil
	},
}

func currentStackRef() string {
	// Allow override with a persistent flag or other non working-directory state.
	return cmdutil.MustGetwd()
}

// Supplies the reserved variable "currentStack" and exits if there is no
// current stack. The supplied query must have a pointer field named
// `Stack` tagged with `graphql:"stackByRef(ref: $currentStack)"`.
func mustQueryStack(ctx context.Context, q any, vars map[string]any) {
	vars = jsonutil.Merge(map[string]any{
		"currentStack": currentStackRef(),
	}, vars)
	if err := api.Query(ctx, svc, q, vars); err != nil {
		cmdutil.Fatalf("query error: %w", err)
	}
	if reflect.ValueOf(q).Elem().FieldByName("Stack").IsNil() {
		cmdutil.Fatalf("no current stack")
	}
}
