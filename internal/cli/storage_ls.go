package cli

import (
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	storageCmd.AddCommand(storageLSCmd)
}

var storageLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "List stores",
	Long:  "List stores in the current stack.", // TODO: In the given scope.
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var q struct {
			Stack *struct {
				Stores []struct {
					Type    string
					Name    string
					SizeMiB *float64
				}
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, nil)
		stores := q.Stack.Stores
		w := cmdutil.NewTableWriter("TYPE", "COMPONENT", "SIZE")
		for _, store := range stores {
			size := ""
			if store.SizeMiB != nil {
				size = cmdutil.FormatBytes(uint64(*store.SizeMiB) * 1024 * 1024)
			}
			w.WriteRow(store.Type, store.Name, size)
		}
		w.Flush()
		return nil
	},
}
