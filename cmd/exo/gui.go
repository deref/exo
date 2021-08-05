package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(guiCmd)
}

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Open exo gui in a web browser",
	Long: `Opens the exo gui in a web browser.

If the current directory is part of a workspace, navigates to it.`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ensureDaemon()
		cl := newClient()

		cwd := cmdutil.MustGetwd()
		output, err := cl.Kernel().FindWorkspace(ctx, &api.FindWorkspaceInput{
			Path: cwd,
		})
		if err != nil {
			return fmt.Errorf("finding workspace: %w", err)
		}

		var endpoint string
		if output.ID == nil {
			endpoint = runState.URL + "/#/new-workspace?root=" + url.QueryEscape(cwd)
		} else {
			endpoint = guiWorkspaceURL(*output.ID)
		}

		fmt.Println("Opening GUI:", endpoint)

		browser.Stdout = os.Stderr
		return browser.OpenURL(endpoint)
	},
}

func guiWorkspaceURL(id string) string {
	return runState.URL + "/#/workspaces/" + url.PathEscape(id)
}
