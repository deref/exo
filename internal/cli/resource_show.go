package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceShowCmd)
}

var resourceShowCmd = &cobra.Command{
	Use:   "show <ref>",
	Short: "Show a resource",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ref := args[0]
		var q struct {
			Resource *struct {
				ID    string
				Type  string
				IRI   *string
				Model scalars.JSONObject
			} `graphql:"resourceByRef(ref: $ref)"`
		}
		if err := api.Query(ctx, svc, &q, map[string]any{
			"ref": ref,
		}); err != nil {
			return err
		}
		if q.Resource != nil {
			resource := *q.Resource
			cmdutil.PrintCueStruct(map[string]any{
				"id":    resource.ID,
				"type":  resource.Type,
				"iri":   resource.IRI,
				"model": (map[string]any)(resource.Model),
			})
		}
		return nil
	},
}
