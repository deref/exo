package cli

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(secretsCmd)
}

var secretsCmd = &cobra.Command{
	Use:  "secrets",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOrEnsureServer()
		ctx := newContext()
		cl := newClient()
		kernel := cl.Kernel()

		authResult, err := kernel.AuthEsv(ctx, &api.AuthEsvInput{})
		if err != nil {
			return fmt.Errorf("getting auth url: %w", err)
		}

		// This link is not opened automatically because it's single use only. That
		// means if we open it in the wrong browser it becomes worthless.
		fmt.Println("Open the following URL to authenticate:")
		fmt.Println(authResult.AuthUrl)
		return nil
	},
}
