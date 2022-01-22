package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobCancelCmd)
}

var jobCancelCmd = &cobra.Command{
	Use:   "cancel <id>",
	Short: "Cancel a job",
	Long:  `Cancel all tasks in a job.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var res struct{}
		return svc.Do(ctx, `
			mutation ($id: String!) {
				cancelJob(id: $id) {
					__typename
				}
			}
		`, map[string]interface{}{
			"id": args[0],
		}, &res)
	},
}
