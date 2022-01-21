package cli

import (
	"fmt"
	"strings"

	"github.com/deref/exo/internal/api"
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
				Model string
			} `graphql:"resourceByRef(ref: $ref)"`
		}
		if err := api.Query(ctx, svc, &q, map[string]interface{}{
			"ref": ref,
		}); err != nil {
			return err
		}
		if q.Resource == nil {
			return nil
		}
		resource := *q.Resource
		fmt.Printf("ID: %s\n", resource.ID)
		fmt.Printf("Type: %s\n", resource.Type)
		if resource.IRI != nil {
			fmt.Printf("IRI: %s\n", *resource.IRI)
		}
		fmt.Printf("Content-Length: %d\n", len(resource.Model))
		fmt.Println()
		fmt.Print(resource.Model)
		if !strings.HasSuffix(resource.Model, "\n") {
			fmt.Println()
		}
		return nil
	},
}
