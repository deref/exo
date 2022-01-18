package cli

import (
	"context"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobMutateCmd)
}

var jobMutateCmd = &cobra.Command{
	Use:   "mutate <mutation> <variable...>",
	Short: "Run a mutation as a job.",
	Long: `Run a mutation as a job, waiting or not as per the global --async flag.

Variables form a JSON object and top-level key/value pairs are expressed in one
of two forms:

variable=string
variable:=raw

Where variable and string are unquoted JSON strings and raw is an encoded JSON
value.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		mutation := args[0]
		vars, err := cmdutil.ArgsToJsonObject(args[1:])
		if err != nil {
			return err
		}
		return sendMutation(ctx, mutation, vars)
	},
}

// TODO: Most usages of watchJob should prefer this.
// TODO: Many usages of this, should probably have codegen'd, type-safe interfaces.
func sendMutation(ctx context.Context, mutation string, vars map[string]interface{}) error {
	jobID, err := svc.Enqueue(ctx, mutation, vars)
	if err != nil {
		return err
	}
	if rootPersistentFlags.Async {
		return nil
	}
	return watchJob(ctx, jobID)
}
