package cli

import (
	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:    "logout",
	Args:   cobra.NoArgs,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl := newClient()
		kernel := cl.Kernel()

		_, err := kernel.UnauthEsv(ctx, &api.UnauthEsvInput{})
		return err
	},
}
