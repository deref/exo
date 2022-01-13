package cli

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"

	gqlclient "github.com/deref/exo/internal/client"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/shurcooL/graphql"
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
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Stack *struct {
				ID   string
				Name string
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, cl, &q, nil)

		stack := q.Stack
		if stack == nil {
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", stack.ID)
		_ = w.Flush()
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
func mustQueryStack(ctx context.Context, cl *gqlclient.Client, q interface{}, vars map[string]interface{}) {
	vars = jsonutil.Merge(map[string]interface{}{
		"currentStack": graphql.String(currentStackRef()),
	}, vars)
	if err := cl.Query(ctx, q, vars); err != nil {
		cmdutil.Fatalf("query error: %w", err)
	}
	if reflect.ValueOf(q).Elem().FieldByName("Stack").IsNil() {
		cmdutil.Fatalf("no current stack")
	}
}
