package main

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit <component>",
	Short: "Edit component spec",
	Long:  "Edit component spec using your preferred editor.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()

		workspace := requireCurrentWorkspace(ctx, cl)
		description, err := workspace.DescribeComponents(ctx, &api.DescribeComponentsInput{
			Refs: []string{componentRef},
		})
		if err != nil {
			return fmt.Errorf("describing components: %w", err)
		}
		if len(description.Components) != 1 {
			return fmt.Errorf("no such component: %q", componentRef)
		}
		component := description.Components[0]

		oldSpec := component.Spec
		newSpec, err := term.EditString("spec.*", oldSpec) // TODO: Correct file extension.

		_, err = workspace.UpdateComponent(ctx, &api.UpdateComponentInput{
			Ref:  component.ID,
			Spec: newSpec,
		})
		// TODO: This should handle a job id for the update step.
		return err
	},
}
