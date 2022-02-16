package cli

import (
	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&lsFlags.All, "all", "a", false, "Show disposed components")
}

var lsFlags struct {
	All bool
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
					ID       string
					Name     string
					Type     string
					Disposed *scalars.Instant
				} `graphql:"components(all: $all)"`
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, map[string]interface{}{
			"all": lsFlags.All,
		})
		var w *cmdutil.TableWriter
		if lsFlags.All {
			w = cmdutil.NewTableWriter("NAME", "ID", "TYPE", "DISPOSED")
			for _, component := range q.Stack.Components {
				w.WriteRow(component.Name, component.ID, component.Type, component.Disposed.String())
			}
		} else {
			w = cmdutil.NewTableWriter("NAME", "ID", "TYPE")
			for _, component := range q.Stack.Components {
				w.WriteRow(component.Name, component.ID, component.Type)
			}
		}
		w.Flush()
		return nil
	},
}
