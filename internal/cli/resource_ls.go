package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceLSCmd)

	resourceLSCmd.AddCommand(makeHelpSubcmd())
}

var resourceLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "List resources",
	Long:  `Lists resources for the current stack.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Stack *struct {
				ID        string
				Resources []struct {
					IRI string
				}
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, cl, &q, nil)
		for _, resource := range q.Stack.Resources {
			fmt.Println(resource.IRI)
		}
		return nil
	},
}
