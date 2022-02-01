package cli

import (
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show <ref>",
	Short: "Show a component",
	Long:  `Show a component's effective configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ref := args[0]
		var q struct {
			Component *struct {
				Configuration string
			} `graphql:"componentByRef(ref: $ref, stack: $stack)"`
		}
		if err := api.Query(ctx, svc, &q, map[string]interface{}{
			"ref":   ref,
			"stack": currentStackRef(),
		}); err != nil {
			return err
		}
		if q.Component == nil {
			return fmt.Errorf("no such component: %q", ref)
		}
		cmdutil.Show(q.Component.Configuration)
		return nil
	},
}
