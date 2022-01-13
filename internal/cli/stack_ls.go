package cli

import (
	"fmt"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackLSCmd)
	stackLSCmd.Flags().BoolVarP(&stackLSFlags.All, "all", "a", false, "List all stacks")
}

var stackLSFlags struct {
	All bool
}

var stackLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists stacks",
	Long: `Lists stacks.

Unless --all is set, scopes stacks to the current project.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		type stackFragment struct {
			ID      string
			Name    string
			Project *struct {
				ID string
			}
		}

		var stacks []stackFragment
		if stackLSFlags.All {
			var q struct {
				Stacks []stackFragment `graphql:"allStacks"`
			}
			if err := cl.Query(ctx, &q, nil); err != nil {
				return fmt.Errorf("querying: %w", err)
			}
			stacks = q.Stacks
		} else {
			var q struct {
				Workspace *struct {
					Project *struct {
						Stacks []stackFragment
					}
				} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
			}
			mustQueryWorkspace(ctx, cl, &q, nil)
			if q.Workspace.Project == nil {
				return fmt.Errorf("no current project")
			}
			stacks = q.Workspace.Project.Stacks
		}

		var w *cmdutil.TableWriter
		if stackLSFlags.All {
			w = cmdutil.NewTableWriter("ID", "NAME", "PROJECT")
			for _, stack := range stacks {
				project := ""
				if stack.Project != nil {
					project = stack.Project.ID
				}
				w.WriteRow(stack.ID, stack.Name, project)
			}
		} else {
			w = cmdutil.NewTableWriter("ID", "NAME")
			for _, stack := range stacks {
				w.WriteRow(stack.ID, stack.Name)
			}
		}
		w.Flush()
		return nil
	},
}
