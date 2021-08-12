package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/deref/exo/internal/core/api"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobWatchCmd)
}

var jobWatchCmd = &cobra.Command{
	Use:   "watch <job-id>",
	Short: "Lists a job's tasks until completion",
	Long:  `Lists a job's tasks as a tree. Rerenders until the job has finished running.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		kernel := cl.Kernel()

		return watchJob(ctx, kernel, args[0])
	},
}

func watchJob(ctx context.Context, kernel api.Kernel, jobID string) error {
	out := os.Stdout

	// Print link to job in GUI.
	routes := newGUIRoutes()
	fmt.Fprintln(out, "Job URL:", routes.JobURL(jobID))

	// Refresh rate starts fast, in case the job completes fast, but will
	// slow over time to minimize overhead and UI flicker.
	delay := 5.0

	w := &lineCountingWriter{
		Underlying: out,
	}

	jp := &jobPrinter{}
	jp.Spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	var job api.TaskDescription
loop:
	for {
		clearLines(w.LineCount)
		w.LineCount = 0

		output, err := kernel.DescribeTasks(ctx, &api.DescribeTasksInput{
			JobIDs: []string{jobID},
		})
		if err != nil {
			return fmt.Errorf("describing tasks: %w", err)
		}
		for _, task := range output.Tasks {
			if task.ID == jobID {
				job = task
				break
			}
		}
		if job.ID == "" {
			return fmt.Errorf("no such job: %q", jobID)
		}

		jp.printTree(w, output.Tasks)

		if job.Finished != nil {
			break
		}
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(time.Duration(delay) * time.Millisecond):
			if delay < 100 {
				delay *= 1.3
			}
		}
		jp.Iteration++
	}

	switch job.Status {
	case taskapi.StatusFailure:
		if job.Message == "" {
			return errors.New("job failed")
		} else {
			return fmt.Errorf("job failure: %s", job.Message)
		}
	case taskapi.StatusSuccess:
		return nil
	default:
		return fmt.Errorf("unexpected job status: %q", job.Status)
	}
}

const esc = 27

var clearLine = fmt.Sprintf("%c[%dA%c[2K", esc, 1, esc)

func clearLines(n int) {
	_, _ = fmt.Fprint(os.Stdout, strings.Repeat(clearLine, n))
}

type lineCountingWriter struct {
	Underlying io.Writer
	LineCount  int
}

func (w *lineCountingWriter) Write(bs []byte) (n int, err error) {
	for _, c := range bs {
		if c == '\n' {
			w.LineCount++
		}
	}
	return w.Underlying.Write(bs)
}
