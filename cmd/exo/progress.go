package main

import (
	"context"
	"fmt"
	"time"

	"github.com/deref/exo/core/api"
	taskapi "github.com/deref/exo/task/api"
	"github.com/deref/exo/util/cmdutil"
)

func showProgress(ctx context.Context, kernel api.Kernel, jobID string) {
	// Refresh rate starts fast, in case the job completes fast, but will
	// slow over time to minimize overhead and UI flicker.
	delay := 10.0

	var job api.TaskDescription
loop:
	for {
		output, err := kernel.DescribeTasks(ctx, &api.DescribeTasksInput{
			JobIDs: []string{jobID},
		})
		if err != nil {
			cmdutil.Fatalf("describing tasks: %w", err)
		}

		for _, task := range output.Tasks {
			fmt.Println(task.ID, task.Name, task.Status, task.Message)
			if task.ID == jobID {
				job = task
			}
		}
		if job.Finished != nil {
			break
		}
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(time.Duration(delay) * time.Millisecond):
			// Refresh at least twice per second.
			if delay < 500 {
				delay *= 1.5
			}
		}
		fmt.Println()
	}

	switch job.Status {
	case taskapi.StatusFailure:
		if job.Message == "" {
			cmdutil.Fatalf("job failed")
		} else {
			cmdutil.Fatalf("job failed: %s", job.Message)
		}
	case taskapi.StatusSuccess:
		// No-op.
	default:
		cmdutil.Fatalf("job status:", job.Status)
	}
}
