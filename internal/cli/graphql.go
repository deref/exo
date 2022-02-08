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

		newRes := func() interface{} {
			var res interface{}
			return &res
		}
		sub := svc.Subscribe(ctx, newRes, doc, vars)
		defer sub.Stop()
		for event := range sub.Events() {
			if err := enc.Encode(event); err != nil {
				return fmt.Errorf("encoding result: %w", err)
			}
		}
		return sub.Err()
	},
}
