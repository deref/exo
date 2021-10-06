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
	guiCmd.Flags().BoolVar(&guiFlags.Print, "print", false, "prints the GUI URL to stdout without launching a browser")
}

var guiFlags struct {
	Print bool
}

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Open exo gui in a web browser",
	Long: `Opens the exo gui in a web browser.

If the current directory is part of a workspace, navigates to it.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		checkOrEnsureServer()
		cl := newClient()

		cwd := cmdutil.MustGetwd()
		output, err := cl.Kernel().ResolveWorkspace(ctx, &api.ResolveWorkspaceInput{
			Ref: cwd,
		})
		if err != nil {
			return fmt.Errorf("resolving workspace: %w", err)
		}

		routes := newGUIRoutes()
		var endpoint string
		if output.ID == nil {
			endpoint = routes.NewWorkspaceURL(cwd)
		} else {
			endpoint = routes.WorkspaceURL(*output.ID)
		}

		// Add a token to auth.
		u, err := url.Parse(endpoint)
		if err != nil {
			return fmt.Errorf("parsing endpoint: %w", err)
		}
		query := u.Query()
		query.Add("token", mustGetToken())
		u.RawQuery = query.Encode()
		endpoint = u.String()

		if guiFlags.Print {
			fmt.Println(endpoint)
			return nil
		}

		fmt.Fprintf(os.Stderr, "Opening GUI: %s\n", endpoint)

		browser.Stdout = os.Stderr
		return browser.OpenURL(endpoint)
	},
}
