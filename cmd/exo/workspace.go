package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/core/client"
	"github.com/deref/exo/internal/util/cmdutil"
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
		ctx := newContext()
		checkOrEnsureServer()

		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.Describe(ctx, &api.DescribeInput{})
		if err != nil {
			cmdutil.Fatalf("describing workspace: %w", err)
		}
		desc := output.Description
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", desc.ID)
		_, _ = fmt.Fprintf(w, "path:\t%s\n", desc.Root)
		_, _ = fmt.Fprintf(w, "display-name:\t%s\n", desc.DisplayName)
		_ = w.Flush()
		return nil
	},
}

func requireCurrentWorkspace(ctx context.Context, cl *client.Root) *client.Workspace {
	workspace := mustResolveCurrentWorkspace(ctx, cl)
	if workspace == nil {
		cmdutil.Fatalf("no workspace for current directory")
	}
	return workspace
}

func mustResolveCurrentWorkspace(ctx context.Context, cl *client.Root) *client.Workspace {
	workspace, err := resolveCurrentWorkspace(ctx, cl)
	if err != nil {
		cmdutil.Fatalf("error resolving workspace: %v", err)
	}
	return workspace
}

func resolveCurrentWorkspace(ctx context.Context, cl *client.Root) (*client.Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getwd: %w", err)
	}
	return resolveWorkspace(ctx, cl, cwd)
}

func resolveWorkspace(ctx context.Context, cl *client.Root, ref string) (*client.Workspace, error) {
	output, err := cl.Kernel().ResolveWorkspace(ctx, &api.ResolveWorkspaceInput{
		Ref: ref,
	})
	if err != nil {
		return nil, err
	}
	var workspace *client.Workspace
	if output.ID != nil {
		workspace = cl.GetWorkspace(*output.ID)
	}
	return workspace, nil
}
