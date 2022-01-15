package cli

import (
	"context"
	"fmt"

	"github.com/alessio/shellescape"
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
		envv, err := getEnvv(ctx)
		if err != nil {
			return err
		}
		for _, kvp := range envv {
			fmt.Println(kvp)
		}
		return nil
	},
}

func getEnvv(ctx context.Context) ([]string, error) {
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

	envVars := q.Workspace.Environment.Variables
	envv := make([]string, len(envVars))
	for i, envVar := range envVars {
		envv[i] = fmt.Sprintf("%s=%s", envVar.Name, shellescape.Quote(envVar.Value))
	}
	return envv, nil
}
