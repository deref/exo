package cli

import (
	"fmt"
	"net/url"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:    "login",
	Args:   cobra.NoArgs,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOrEnsureServer()
		serverURL := effectiveServerURL()
		uri, err := url.Parse(serverURL)
		if err != nil {
			return fmt.Errorf("parsing server url %q: %w", serverURL, err)
		}

		uri.Fragment = "/auth-esv"
		authURL, err := addAuthTokenToURL(uri.String())
		if err != nil {
			return fmt.Errorf("adding auth token to url %q: %w", uri.String(), err)
		}

		fmt.Println("Opening " + authURL)
		if err := browser.OpenURL(authURL); err != nil {
			fmt.Println("Failed to open a browser. Please open the above link manually.")
		}
		return nil
	},
}
