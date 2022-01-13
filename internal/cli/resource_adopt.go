package cli

import (
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceAdoptCmd)
}

var resourceAdoptCmd = &cobra.Command{
	Use:   "adopt <iri>",
	Short: "Adopt a resource",
	Long:  "Adopt an existing resource.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		iri := args[0]

		var m struct {
			Resource struct {
				IRI string
			} `graphql:"adoptResource(iri: $iri)"`
		}
		return cl.Mutate(ctx, &m, map[string]interface{}{
			"iri": graphql.String(iri),
		})
	},
}
