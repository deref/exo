package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().StringVar(&workerFlags.Job, "job", "", "job to work")
}

var workerFlags struct {
	Job string
}

var workerCmd = &cobra.Command{
	Hidden: true,
	Use:    "worker",
	Short:  "Run a task worker",
	Long: `Run a task worker.

If --job is specified, only works tasks for that specific job and terminates
when the job terminates.
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if !isPeerMode() {
			return errors.New("worker command only available in peer mode")
		}

		worker := &api.Worker{
			Service: svc,
			ID:      fmt.Sprintf("peer:%d:worker", os.Getpid()),
			JobID:   workerFlags.Job,
		}
		return worker.Run(ctx, ctx)
	},
}
