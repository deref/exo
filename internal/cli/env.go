package cli

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/util/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
	// TODO: --scope, --component=, etc.
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Long:  `Prints the workspace's environment variables in .env format.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		showEnvironment(getWorkspaceEnvironment(ctx))
		return nil
	},
}

func getWorkspaceEnvironment(ctx context.Context) environmentFragment {
	var q struct {
		Workspace *struct {
			Environment environmentFragment
		} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
	}
	mustQueryWorkspace(ctx, &q, nil)
	return q.Workspace.Environment
}

func showEnvironment(env environmentFragment) {
	for _, v := range env.Variables {
		fmt.Println(osutil.FormatDotEnvEntry(v.Name, v.Value))
	}
}

type environmentFragment struct {
	Variables []struct {
		Name  string
		Value string
	}
}

func environmentFragmentToMap(fragment environmentFragment) map[string]string {
	variables := fragment.Variables
	m := make(map[string]string, len(variables))
	for _, variable := range variables {
		m[variable.Name] = variable.Value
	}
	return m
}
