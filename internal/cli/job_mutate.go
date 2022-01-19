package cli

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/peer"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
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
	Hidden: true,
	Args:   cobra.MinimumNArgs(1),
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
	jobID, err := api.Enqueue(ctx, svc, mutation, vars)
	if err != nil {
		return err
	}
	if rootPersistentFlags.Async {
		fmt.Println("job:", jobID)
		return nil
	}
	var eg errgroup.Group

	// When the CLI is in peer mode, there is generally no worker pool.  When
	// awaiting job completion, do the work as part of this CLI invocation.
	if p, ok := svc.(*peer.Peer); ok {
		eg.Go(func() error {
			return peer.WorkTask(ctx, p, jobID)
		})
	}

	eg.Go(func() error {
		return watchJob(ctx, jobID)
	})

	return eg.Wait()
}
