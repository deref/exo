package cli

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().BoolVarP(&treeFlags.All, "all", "a", false, "Show disposed components")
}

var treeFlags struct {
	All bool
}

var treeCmd = &cobra.Command{
	Hidden: true, // Most users don't need to care about subcomponents.
	Use:    "tree",
	Short:  "Print component tree",
	Long:   `Prints a tree of all components in the current stack.`,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		type componentFragment struct {
			ParentID *string
			ID       string
			Name     string
			Type     string
			Disposed *scalars.Instant
		}
		var q struct {
			Stack *struct {
				Components []componentFragment `graphql:"components(all: $all, recursive: true)"`
			} `graphql:"stackByRef(ref: $currentStack)"`
		}
		mustQueryStack(ctx, &q, map[string]interface{}{
			"all": treeFlags.All,
		})

		tb := term.NewTreeBuilder()
		for _, component := range q.Stack.Components {
			node := &term.TreeNode{
				ID:     component.ID,
				Label:  fmt.Sprintf("%s (%s)", component.Name, component.Type),
				Suffix: component.ID,
			}
			if component.Disposed != nil {
				node.Label += " DISPOSED"
			}
			if component.ParentID != nil {
				node.ParentID = *component.ParentID
			}
			tb.AddNode(node)
		}
		trees := tb.Build()
		for _, tree := range trees {
			term.PrintTree(os.Stdout, tree)
		}
		return nil
	},
}
