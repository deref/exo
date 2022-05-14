package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackShowCmd)
}

var stackShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show stack details",
	Long:  `Show details about a stack.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var q struct {
			Stack *struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Cluster     struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"cluster"`
				Project *struct {
					ID string `json:"id"`
				} `json:"project"`
				Disposed *api.Instant `json:"disposed"`
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, nil)
		cmdutil.PrintCueStruct(q.Stack)
		return nil
	},
}
