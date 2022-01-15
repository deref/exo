package cli

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"

	gqlclient "github.com/deref/exo/internal/client"
	"github.com/deref/exo/internal/core/api"
	joshclient "github.com/deref/exo/internal/core/client"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workspaceCmd)

	workspaceCmd.AddCommand(makeHelpSubcmd())
}

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Create, inspect, and modify workspaces",
	Long: `Contains subcommands for operating on workspaces.

If no subcommand is given, describes the current workspace.`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return nil
		}
		ctx := cmd.Context()

		var q struct {
			Workspace *struct {
				ID   string
				Root string
			} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
		}
		mustQueryWorkspace(ctx, client, &q, nil)
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", q.Workspace.ID)
		_, _ = fmt.Fprintf(w, "path:\t%s\n", q.Workspace.Root)
		_ = w.Flush()
		return nil
	},
}

func requireCurrentWorkspace(ctx context.Context, cl *joshclient.Root) *joshclient.Workspace {
	workspace := mustResolveCurrentWorkspace(ctx, cl)
	if workspace == nil {
		cmdutil.Fatalf("no workspace for current directory")
	}
	return workspace
}

func mustResolveCurrentWorkspace(ctx context.Context, cl *joshclient.Root) *joshclient.Workspace {
	workspace, err := resolveCurrentWorkspace(ctx, cl)
	if err != nil {
		cmdutil.Fatalf("error resolving workspace: %v", err)
	}
	return workspace
}

func currentWorkspaceRef() string {
	// Allow override with a persistent flag or other non working-directory state.
	return cmdutil.MustGetwd()
}

func resolveCurrentWorkspace(ctx context.Context, cl *joshclient.Root) (*joshclient.Workspace, error) {
	return resolveWorkspace(ctx, cl, currentWorkspaceRef())
}

func resolveWorkspace(ctx context.Context, cl *joshclient.Root, ref string) (*joshclient.Workspace, error) {
	output, err := cl.Kernel().ResolveWorkspace(ctx, &api.ResolveWorkspaceInput{
		Ref: ref,
	})
	if err != nil {
		return nil, err
	}
	var workspace *joshclient.Workspace
	if output.ID != nil {
		workspace = cl.GetWorkspace(*output.ID)
	}
	return workspace, nil
}

// Supplies the reserved variable "currentWorkspace" and exits if there is no
// current workspace. The supplied query must have a pointer field named
// `Workspace` tagged with `graphql:"workspaceByRef(ref: $currentWorkspace)"`.
func mustQueryWorkspace(ctx context.Context, cl *gqlclient.Client, q interface{}, vars map[string]interface{}) {
	vars = jsonutil.Merge(map[string]interface{}{
		"currentWorkspace": currentWorkspaceRef(),
	}, vars)
	if err := cl.Query(ctx, q, vars); err != nil {
		cmdutil.Fatalf("query error: %w", err)
	}
	if reflect.ValueOf(q).Elem().FieldByName("Workspace").IsNil() {
		cmdutil.Fatalf("no current workspace")
	}
}
