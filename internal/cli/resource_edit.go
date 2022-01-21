package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceEditCmd)
}

var resourceEditCmd = &cobra.Command{
	Use:   "edit <ref>",
	Short: "Edit a resource model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ref := args[0]
		var q struct {
			Resource *struct {
				Model string
			} `graphql:"resourceByRef(ref: $ref)"`
		}
		if err := api.Query(ctx, svc, &q, map[string]interface{}{
			"ref": ref,
		}); err != nil {
			return fmt.Errorf("querying: %w", err)
		}
		if q.Resource == nil {
			return errors.New("no such resource")
		}
		newModel, err := term.EditString("resource.*.txt", q.Resource.Model)
		if strings.TrimSpace(newModel) == "" {
			return errors.New("aborting due to empty model")
		}
		if err != nil {
			return fmt.Errorf("editing: %w", err)
		}
		return sendMutation(ctx, "updateResource", map[string]interface{}{
			"ref":   ref,
			"model": newModel,
		})
	},
}
