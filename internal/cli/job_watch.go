package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/deref/exo/internal/api"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobWatchCmd)
}

var jobWatchCmd = &cobra.Command{
	Use:   "watch <job-id>",
	Short: "Watch a job's progress",
	Long: `Tails the events from all tasks in a job and continuously renders a task
tree until the job has finished running.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return watchJob(ctx, args[0])
	},
}

type jobEventFragment struct {
	Type    string
	Message string
	Job     struct {
		URL      string
		Tasks    []taskFragment
		RootTask taskFragment
	}
}

func watchJob(ctx context.Context, jobID string) error {
	out := os.Stdout

	var res struct {
		Event jobEventFragment `graphql:"watchJob(id: $id)"`
	}
	sub := api.Subscribe(ctx, svc, &res, map[string]interface{}{
		"id": jobID,
	})
	defer sub.Stop()

	w := &lineCountingWriter{
		Underlying: out,
	}

	jp := &jobPrinter{}
	jp.Spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	var prev jobEventFragment

	// Periodically tick to keep spinner animation lively, even when
	// there are no events. Start ticking only after the first event.
	var tickC <-chan time.Time

watching:
	for {
		var event jobEventFragment
		select {
		case more, _ := <-sub.C():
			if !more {
				break watching
			}
			event = res.Event
		case <-tickC:
			event = jobEventFragment{
				Type: "Tick",
				Job:  prev.Job,
			}
		}
		job := event.Job

		clearLines(w.LineCount)
		w.LineCount = 0

		switch event.Type {
		case "Tick":
			jp.Iteration++

		case "JobWatched":
			fmt.Fprintln(out, "Job URL:", job.URL)
			if tickC != nil {
				return errors.New("already received JobWatched event")
			}
			ticker := time.NewTicker(time.Second / time.Duration(len(jp.Spinner)))
			tickC = ticker.C
			defer ticker.Stop()

		case "JobUpdated", "TaskFinished":
			// No-op.

		default:
			// XXX print with colored header etc a la logs.
			fmt.Fprintln(out, event.Message)
		}

		jp.printTree(w, job.Tasks)

		if job.RootTask.Finished != nil {
			sub.Stop()
		}

		prev = event
	}
	if sub.Err() != nil {
		return fmt.Errorf("subscription error: %w", sub.Err())
	}

	root := prev.Job.RootTask
	switch root.Status {
	case taskapi.StatusFailure:
		if root.Message == "" {
			return errors.New("job failed")
		} else {
			return fmt.Errorf("job failure: %s", root.Message)
		}
	case taskapi.StatusSuccess:
		return nil
	default:
		return fmt.Errorf("unexpected job status: %q", root.Status)
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
