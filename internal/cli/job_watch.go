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

type eventFragment struct {
	Message string
}

func watchJob(ctx context.Context, jobID string) error {
	out := os.Stdout

	var res struct {
		Output struct {
			Job struct {
				URL   string
				Tasks []taskFragment
			}
			Event *eventFragment
		} `graphql:"watchJob(id: $id)"`
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

	first := true
	var root taskFragment

	// Periodically tick to keep spinner animation lively, even when
	// there are no events. Start ticking only after the first event.
	var tickC <-chan time.Time

watching:
	for {
		fmt.Println("watch loop")
		var event *eventFragment
		select {
		case more, _ := <-sub.C():
			fmt.Println("more?", more)
			fmt.Printf("res= %#v\n", res)
			if !more {
				break watching
			}
			event = res.Output.Event
		case <-tickC:
			fmt.Println("tickC")
			event = nil
		}
		job := res.Output.Job

		if first {
			fmt.Fprintln(out, "Job URL:", job.URL)
			first = false
			ticker := time.NewTicker(time.Second / time.Duration(len(jp.Spinner)))
			tickC = ticker.C
			defer ticker.Stop()
		}

		//clearLines(w.LineCount)
		w.LineCount = 0

		if event != nil {
			// XXX print with colored header etc a la logs.
			fmt.Fprintln(out, event.Message)
		}

		for _, task := range job.Tasks {
			if task.ParentID == nil {
				root = task
				break
			}
		}
		if root.ID == "" {
			return errors.New("root task missing")
		}

		jp.printTree(w, job.Tasks)

		if root.Finished != nil {
			break
		}
		// XXX smooth out animation
		jp.Iteration++
	}

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
