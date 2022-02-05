package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(graphqlCmd)
}

var graphqlCmd = &cobra.Command{
	Use:    "graphql <doc> [variables...]",
	Hidden: true,
	Short:  "Executes a graphql operation",
	Long: `Executes a graphql operation.

Arguments are specified as JSON with the same syntax as 'exo json'.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		doc := args[0]
		vars, err := cmdutil.ArgsToJsonObject(args[1:])
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		var res interface{}
		sub := svc.Subscribe(ctx, &res, doc, vars)
		defer sub.Stop()
		for range sub.C() {
			if err := enc.Encode(res); err != nil {
				return fmt.Errorf("encoding result: %w", err)
			}
		}
		return sub.Err()
	},
}
