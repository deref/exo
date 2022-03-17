package cli

import (
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&showFlags.Recursive, "recursive", false, "include subcomponents")
	showCmd.Flags().BoolVar(&showFlags.Final, "final", false, "fully evaluate to a concrete value")
}

var showFlags struct {
	Recursive bool
	Final     bool
}

// TODO: Generalize to be able to conveniently show any type of object.
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
				Configuration string `graphql:"configuration(recursive: $recursive, final: $final)"`
			} `graphql:"componentByRef(ref: $ref, stack: $stack)"`
		}
		if err := api.Query(ctx, svc, &q, map[string]any{
			"ref":       ref,
			"stack":     currentStackRef(),
			"recursive": showFlags.Recursive,
			"final":     showFlags.Final,
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
