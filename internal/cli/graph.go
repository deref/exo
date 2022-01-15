package cli

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(graphCmd)
}

var graphCmd = &cobra.Command{
	Use:    "graph",
	Short:  "Render a component graph",
	Long:   "Renders a component graph in dot format.",
	Args:   cobra.NoArgs,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl := newClient()

		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.RenderDependencies(ctx, &api.RenderDependenciesInput{})
		if err != nil {
			return err
		}

		fmt.Print(output.Dot)
		return nil
	},
}
