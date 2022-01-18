package cli

import (
	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm [ref ...]",
	Short: "Remove components",
	Long:  "Remove components.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		ctx := cmd.Context()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.DeleteComponents(ctx, &api.DeleteComponentsInput{
			Refs: args,
		})
		if err != nil {
			return err
		}
		return watchJob(ctx, output.JobID)
	},
}
