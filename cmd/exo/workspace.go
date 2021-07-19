package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/exod/api"
	"github.com/deref/exo/exod/client"
	"github.com/deref/exo/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workspaceCmd)

	workspaceCmd.AddCommand(helpSubcmd)
}

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Create, inspect, and modify workspaces.",
	Long: `Contains subcommands for operating on workspaces.

If no subcommand is given, describes the current workspace.`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return nil
		}
		ctx := newContext()
		ensureDeamon()

		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		output, err := workspace.Describe(ctx, &api.DescribeInput{})
		if err != nil {
			cmdutil.Fatalf("describing workspace: %w", err)
		}
		desc := output.Description
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", desc.ID)
		_, _ = fmt.Fprintf(w, "path:\t%s\n", desc.Root)
		_ = w.Flush()
		return nil
	},
}

func requireWorkspace(ctx context.Context, root *client.Root) api.Workspace {
	workspace := mustFindWorkspace(ctx, root)
	if workspace == nil {
		cmdutil.Fatalf("no workspace for current directory")
	}
	return workspace
}

func mustFindWorkspace(ctx context.Context, root *client.Root) api.Workspace {
	workspace, err := findWorkspace(ctx, root)
	if err != nil {
		cmdutil.Fatal(err)
	}
	return workspace
}

func findWorkspace(ctx context.Context, root *client.Root) (api.Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getwd: %w", err)
	}
	output, err := root.Kernel().FindWorkspace(ctx, &api.FindWorkspaceInput{
		Path: cwd,
	})
	if err != nil {
		return nil, fmt.Errorf("finding workspace: %w", err)
	}
	var workspace api.Workspace
	if output.ID != nil {
		workspace = root.GetWorkspace(*output.ID)
	}
	return workspace, nil
}
