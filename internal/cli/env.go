package cli

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/util/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Long:  `Prints the workspace's environment variables in .env format.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		env := osutil.EnvMapToDotEnv(getEnvMap(ctx))
		for _, kvp := range env {
			fmt.Println(kvp)
		}
		return nil
	},
}

func getEnvMap(ctx context.Context) map[string]string {
	var q struct {
		Workspace *struct {
			Environment struct {
				Variables []struct {
					Name  string
					Value string
				}
			}
		} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
	}
	mustQueryWorkspace(ctx, &q, nil)
	vars := q.Workspace.Environment.Variables
	m := make(map[string]string, len(vars))
	for _, v := range vars {
		m[v.Name] = v.Value
	}
	return m
}
