package cli

import (
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit <component>",
	Short: "Edit component spec",
	Long:  "Edit component spec using your preferred editor.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ref := args[0]

		var q struct {
			Component *struct {
				ID   string
				Spec string
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

		oldSpec := q.Component.Spec
		newSpec, err := term.EditString(ref+".*.cue", oldSpec) // TODO: Correct file extension.
		if err != nil {
			return fmt.Errorf("editing: %w", err)
		}

		return sendMutation(ctx, "updateComponent", map[string]interface{}{
			"ref":     q.Component.ID,
			"newSpec": newSpec,
		})
	},
}
