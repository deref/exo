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
	"github.com/deref/exo/internal/scalars"
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
	Type      string
	Message   string
	Timestamp scalars.Instant
	Job       struct {
		URL      string
		Tasks    []taskFragment
		RootTask taskFragment
	}
	JobID  *string `graphql:"jobID"`
	TaskID *string `graphql:"taskID"`
}

func watchJob(ctx context.Context, jobID string) error {
	out := os.Stdout

	type watchJobSubscription struct {
		Event jobEventFragment `graphql:"watchJob(id: $id)"`
	}
	var res watchJobSubscription
	sub := api.Subscribe(ctx, svc, &res, map[string]interface{}{
		"id": jobID,
	})
	defer sub.Stop()

	w := &EventWriter{
		W: out,
	}
	w.Init()

	lcw := &lineCountingWriter{
		Underlying: out,
	}

	interactive := isInteractive()

	var jp *jobPrinter
	if interactive {
		jp = &jobPrinter{}
		jp.Spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}

	var prev jobEventFragment

	// Periodically tick to keep spinner animation lively, even when
	// there are no events. Start ticking only after the first event.
	var tickC <-chan time.Time

watching:
	for {
		var event jobEventFragment
		select {
		case eventInterface, ok := <-sub.Events():
			if !ok {
				break watching
			}
			event = api.OperationData(eventInterface).(jobEventFragment)
		case <-tickC:
			event = jobEventFragment{
				Type: "Tick",
				Job:  prev.Job,
			}
		}
		job := event.Job

		if interactive {
			clearLines(lcw.LineCount)
			lcw.LineCount = 0
		}

		switch event.Type {
		case "Tick":
			jp.Iteration++

		case "JobWatched":
			fmt.Fprintln(out, "Job URL:", job.URL)
			if tickC != nil {
				return errors.New("already received JobWatched event")
			}
			if interactive {
				ticker := time.NewTicker(time.Second / time.Duration(len(jp.Spinner)))
				tickC = ticker.C
				defer ticker.Stop()
			}

		case "JobUpdated":
			if !interactive {
				// XXX only print occassionally.
				// XXX include progress info.
				sourceID := *event.JobID
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceID, "JobUpdated")
			}

		case "TaskStarted":
			if !interactive {
				sourceID := *event.TaskID
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceID, "Task Started")
			}

		case "TaskFinished":
			if !interactive {
				sourceID := *event.TaskID
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceID, "Task Finished")
			}
			if *event.TaskID == *event.JobID {
				sub.Stop()
			}

		default:
			sourceID := "UNKNOWN" // XXX
			message := event.Type
			if event.Message != "" {
				message = fmt.Sprintf("%s: %s", message, event.Message)
			}
			w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceID, message)
		}

		if interactive {
			jp.printTree(lcw, job.Tasks)
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
