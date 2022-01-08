package cli

import (
	"context"
	"fmt"
	"sort"

	"github.com/alessio/shellescape"
	"github.com/deref/exo/internal/core/api"
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
	checkOrEnsureServer()
	cl := newClient()
	workspace := requireCurrentWorkspace(ctx, cl)
	output, err := workspace.DescribeEnvironment(ctx, &api.DescribeEnvironmentInput{})
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(output.Variables))
	i := 0
	for key := range output.Variables {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	envv := make([]string, i)
	for i, key := range keys {
		value := output.Variables[key]
		envv[i] = fmt.Sprintf("%s=%s", key, shellescape.Quote(value.Value))
	}
	return envv, nil
}
