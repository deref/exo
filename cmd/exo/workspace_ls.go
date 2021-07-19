package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/exod/api"
	"github.com/deref/exo/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	workspaceCmd.AddCommand(workspaceLSCmd)
}

var workspaceLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists workspaces",
	Long:  `Lists workspaces.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDeamon()
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
