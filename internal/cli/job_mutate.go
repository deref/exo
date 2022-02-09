package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	jobCmd.AddCommand(jobMutateCmd)
}

var jobMutateCmd = &cobra.Command{
	Use:   "mutate <mutation> <arguments...>",
	Short: "Run a mutation as a job.",
	Long: `Run a mutation as a job, waiting or not as per the global --async flag.

Arguments are specified as JSON with the same syntax as 'exo json'.`,
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

// TODO: Most usages of watchJob should prefer this or watchOwnJob.
// TODO: Many usages of this, should probably have codegen'd, type-safe interfaces.
func sendMutation(ctx context.Context, mutation string, vars map[string]interface{}) error {
	jobID, err := api.Enqueue(ctx, svc, mutation, vars)
	if err != nil {
		return err
	}
	return watchOwnJob(ctx, jobID)
}

// Watch a job that we just enqueued.
// If so configured, also runs a worker to complete the job.
func watchOwnJob(ctx context.Context, jobID string) error {
	if rootPersistentFlags.Async {
		fmt.Println("job:", jobID)
		return nil
	}
	var eg errgroup.Group

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if workOwnJobs {
		eg.Go(func() error {
			workerID := fmt.Sprintf("cli:%d", os.Getpid())
			err := api.RunWorker(ctx, svc, workerID, &jobID)
			if err != nil {
				cancel()
			}
			return err
		})
	}

	eg.Go(func() error {
		return watchJob(ctx, jobID)
	})

	return eg.Wait()
}
