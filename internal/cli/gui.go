package cli

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/deref/exo/internal/api"
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

func addAuthTokenToURL(uri string) (string, error) {
	// Add a token to auth.
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("parsing endpoint: %w", err)
	}
	query := u.Query()
	query.Add("token", mustGetToken())
	u.RawQuery = query.Encode()
	return u.String(), nil
}

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Open exo gui in a web browser",
	Long: `Opens the exo gui in a web browser.

If the current directory is part of a workspace, navigates to it.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		var q struct {
			Routes struct {
				NewProjectURL string `graphql:"newProjectUrl(workspace: $cwd)"`
			}
			Workspace *struct {
				URL string
			} `graphql:"workspaceByRef(ref: $cwd)"`
		}
		if err := api.Query(ctx, svc, &q, map[string]interface{}{
			"cwd": cmdutil.MustGetwd(),
		}); err != nil {
			return fmt.Errorf("querying: %w", err)
		}

		var endpoint string
		if q.Workspace == nil {
			endpoint = q.Workspace.URL
		} else {
			endpoint = q.Routes.NewProjectURL
		}

		// TODO: Add auth-token server-side?
		endpoint, err := addAuthTokenToURL(endpoint)
		if err != nil {
			return err
		}

		if guiFlags.Print {
			fmt.Println(endpoint)
			return nil
		}

		fmt.Fprintf(os.Stderr, "Opening GUI: %s\n", endpoint)

		browser.Stdout = os.Stderr
		return browser.OpenURL(endpoint)
	},
}
