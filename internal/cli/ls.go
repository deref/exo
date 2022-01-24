package cli

import (
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists components",
	Long:  `Lists components.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var q struct {
			Stack *struct {
				Components []struct {
					ID   string
					Name string
					Type string
				}
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, nil)
		w := cmdutil.NewTableWriter("NAME", "ID", "TYPE")
		for _, component := range q.Stack.Components {
			w.WriteRow(component.Name, component.ID, component.Type)
		}
		w.Flush()
		return nil
	},
}
