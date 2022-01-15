package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceForgetCmd)
}

var resourceForgetCmd = &cobra.Command{
	Use:   "forget <iri>", // TODO: Variadic.
	Short: "Forget a resource",
	Long:  `Stop tracking a resource without attempting to destroy it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		iri := args[0]

		var resp struct{}
		return client.Run(ctx, `
			mutation ($iri: String!) {
				forgetResource(iri: $iri)
			}
		`, &resp, map[string]interface{}{
			"iri": iri,
		})
	},
}
