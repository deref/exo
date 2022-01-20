package cli

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/peer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Hidden: true,
	Use:    "worker",
	Short:  "Run a task worker",
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		p, ok := svc.(*peer.Peer)
		if !ok {
			return errors.New("worker command only available in peer mode")
		}

		workerID := fmt.Sprintf("peer:%d:worker", os.Getpid())
		acquireVars := map[string]interface{}{
			"workerId": workerID,
		}
		for {
			var acquired struct {
				Task struct {
					ID string
				} `graphql:"acquireTask(workerId: $workerId)"`
			}
			if err := api.Mutate(ctx, svc, &acquired, acquireVars); err != nil {
				return fmt.Errorf("acquiring task: %w", err)
			}
			taskID := acquired.Task.ID
			log.Printf("acquired task: %s", taskID)
			err := peer.WorkTask(ctx, p, taskID, workerID)
			if err == nil {
				log.Printf("completed task: %s", taskID)
			} else {
				log.Printf("task %s failure: %v", taskID, err)
			}
		}
	},
}
