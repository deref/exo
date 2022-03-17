package cli

import (
	"errors"
	"fmt"

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
		mustQueryWorkspace(ctx, &q, map[string]any{
			"path": args[0],
		})
		f := q.Workspace.FileSystem.File
		if f == nil {
			return errors.New("no such file or directory")
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
