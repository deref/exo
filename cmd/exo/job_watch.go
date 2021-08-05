package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

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
	Long:  `Lists a job's tasks as a tree. Rerenders until the job has finished running. `,
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
	// Refresh rate starts fast, in case the job completes fast, but will
	// slow over time to minimize overhead and UI flicker.
	delay := 10.0

	w := &lineCountingWriter{
		Underlying: os.Stdout,
	}

	jp := &jobPrinter{}
	jp.Spinner = `/-\\|/-` // TODO: `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`

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
			// Refresh at least twice per second.
			if delay < 500 {
				delay *= 1.5
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

func getTermSize() (w, h int) {
	out, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		return 0, 0
	}
	var size struct {
		rows uint16
		cols uint16
	}
	defer out.Close()
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, out.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&size)))
	return int(size.cols), int(size.rows)
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
