package cli

import (
	"fmt"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

func init() {
	fsCmd.AddCommand(fsLSCmd)
}

var fsLSCmd = &cobra.Command{
	Use:   "ls <path>",
	Short: "List directory contents",
	Long: `List directory contents.

Paths are forward-slash separated and must be absolute.

Returns non-zero exit code if the file or directory does not exist.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Workspace *struct {
				FileSystem *struct {
					File *struct {
						Name        string
						IsDirectory bool
						Children    []struct {
							Name string
						}
					} `graphql:"file(path: $path)"`
				}
			} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
		}
		mustQueryWorkspace(ctx, cl, &q, map[string]interface{}{
			"path": graphql.String(args[0]),
		})
		f := q.Workspace.FileSystem.File
		if f == nil {
			cmdutil.Fatalf("no such file or directory")
		}
		if f.IsDirectory {
			for _, child := range f.Children {
				fmt.Println(child.Name)
			}
		} else {
			fmt.Println(f.Name)
		}
		return nil
	},
}
