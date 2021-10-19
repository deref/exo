package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	workspaceCmd.AddCommand(workspaceLSCmd)
}

var workspaceLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists workspaces",
	Long:  `Lists workspaces.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		output, err := cl.Kernel().DescribeWorkspaces(ctx, &api.DescribeWorkspacesInput{})
		if err != nil {
			cmdutil.Fatalf("describing workspaces: %w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		for _, process := range output.Workspaces {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", process.ID, process.Root)
		}
		_ = w.Flush()
		return nil
	},
}
