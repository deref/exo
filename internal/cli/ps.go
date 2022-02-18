package cli

import (
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List processes",
	Long:  `Describes running processes.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var q struct {
			Stack *struct {
				Processes []struct {
					ID   string
					Name string
				} `graphql:"processes"`
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, nil)
		w := cmdutil.NewTableWriter("NAME", "ID")
		for _, process := range q.Stack.Processes {
			w.WriteRow(process.Name, process.ID)
		}
		w.Flush()
		return nil
	},
}
