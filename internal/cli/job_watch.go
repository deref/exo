package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
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
	Type       string
	Message    string
	Timestamp  scalars.Instant
	SourceType string
	SourceID   string
	Job        *jobEventJobFragment
	JobID      *string
	Task       *struct {
		ID    string
		Error *string
	}
}

type jobEventJobFragment struct {
	URL      string
	Tasks    []taskFragment
	RootTask taskFragment
}

func watchJob(ctx context.Context, jobID string) error {
	// Buffer output to minimize terminal flicker when redrawing tree.
	var buf bytes.Buffer
	flush := func() {
		os.Stdout.Write(buf.Bytes())
		buf.Reset()
	}

	type watchJobSubscription struct {
		Event jobEventFragment `graphql:"watchJob(id: $id, debug: $debug)"`
	}
	var res watchJobSubscription
	sub := api.Subscribe(ctx, svc, &res, map[string]interface{}{
		"id":    jobID,
		"debug": isDebugMode(),
	})
	defer sub.Stop()

	w := &EventWriter{
		W: &buf,
	}
	w.Init()

	lcw := &lineCountingWriter{
		Underlying: &buf,
	}

	interactive := isInteractive()
	verbose := !interactive || isDebugMode()

	var jp *jobPrinter
	if interactive {
		jp = &jobPrinter{}
		jp.Spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}

	var job *jobEventJobFragment

	// Periodically tick to keep spinner animation lively, even when
	// there are no events. Start ticking only after the first event.
	initialized := false
	var tickC <-chan time.Time

watching:
	for {
		flush()

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
			}
		}
		if event.Job != nil {
			job = event.Job
		}
		task := event.Task

		if interactive {
			clearLines(lcw.LineCount)
			lcw.LineCount = 0
		}

		sourceID := event.SourceID
		sourceLabel := sourceID // XXX

		switch event.Type {
		case "Tick":
			jp.Iteration++

		case "JobWatched":
			fmt.Fprintln(&buf, "Job URL:", job.URL)
			if initialized {
				return errors.New("already received JobWatched event")
			}
			if interactive {
				ticker := time.NewTicker(time.Second / time.Duration(len(jp.Spinner)))
				tickC = ticker.C
				defer ticker.Stop()
			}
			initialized = true

		case "JobUpdated":
			if verbose {
				// XXX only print occassionally.
				// XXX include progress info.
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, "JobUpdated")
			}

		case "TaskStarted":
			if verbose {
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, "Task Started")
			}

		case "TaskFinished":
			if verbose {
				var message string
				if task.Error == nil {
					message = "task finished; awaiting children for completion"
				} else {
					message = fmt.Sprintf("task failed: %s", *task.Error)
				}
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
			}

		case "TaskCompleted":
			if verbose {
				var message string
				if task.Error == nil {
					message = "task completed successfully"
				} else {
					message = fmt.Sprintf("task failed: %s", *task.Error)
				}
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
			}
			if task.ID == *event.JobID {
				sub.Stop()
			}

		case "Message":
			w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, event.Message)

		default:
			message := event.Type
			if event.Message != "" {
				message = fmt.Sprintf("%s: %s", message, event.Message)
			}
			w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
		}

		if interactive {
			if initialized {
				// Skip line between log output and task tree.
				fmt.Fprintln(lcw)
			}
			jp.printTree(lcw, job.Tasks)
		}
	}
	if sub.Err() != nil {
		return fmt.Errorf("subscription error: %w", sub.Err())
	}

	if job == nil {
		return errors.New("never received job details")
	}
	root := job.RootTask
	if root.Error != nil {
		return fmt.Errorf("job failure: %s", *root.Error)
	}
	return nil
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
