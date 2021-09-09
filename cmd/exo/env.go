package main

import (
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
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.DescribeEnvironment(ctx, &api.DescribeEnvironmentInput{})
		if err != nil {
			return err
		}

		keys := make([]string, len(output.Variables))
		i := 0
		for key := range output.Variables {
			keys[i] = key
			i++
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := output.Variables[key]
			fmt.Printf("%s=%s\n", key, shellescape.Quote(value))
		}

		return nil
	},
}
