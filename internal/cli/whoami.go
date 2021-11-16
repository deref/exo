package cli

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami [authentication server]",
	Short: "Prints the current user's identity.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOrEnsureServer()
		ctx := newContext()
		cl := newClient()
		kernel := cl.Kernel()

		authServer := "https://secrets.deref.io"
		if len(args) > 0 {
			authServer = args[0]
		}

		user, err := kernel.GetEsvUser(ctx, &api.GetEsvUserInput{VaultURL: authServer})
		if err != nil {
			return err
		}

		if user != nil {
			fmt.Println(user.Email)
		} else {
			fmt.Fprintln(os.Stderr, "Not logged in")
		}

		return nil
	},
}
