package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/deref/exo/internal/core/api"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/util/cmdutil"
)

type taskNode struct {
	Parent   *taskNode
	ID       string
	Created  string
	Status   string
	Name     string
	Message  string
	Children []*taskNode
}

func showProgress(ctx context.Context, kernel api.Kernel, jobID string) {
	// Refresh rate starts fast, in case the job completes fast, but will
	// slow over time to minimize overhead and UI flicker.
	delay := 10.0
	nlines := 0
	iteration := 0

	var job api.TaskDescription
loop:
	for {
		clearLines(nlines)

		output, err := kernel.DescribeTasks(ctx, &api.DescribeTasksInput{
			JobIDs: []string{jobID},
		})
		if err != nil {
			cmdutil.Fatalf("describing tasks: %w", err)
		}

		nlines = 0
		printf := func(format string, v ...interface{}) {
			fmt.Printf(format+"\n", v...)
			nlines++
		}

		nodes := make(map[string]*taskNode)
		getNode := func(id string) *taskNode {
			node := nodes[id]
			if node == nil {
				node = &taskNode{}
				nodes[id] = node
			}
			return node
		}

		for _, task := range output.Tasks {
			child := getNode(task.ID)
			child.ID = task.ID
			child.Created = task.Created
			child.Name = task.Name
			child.Status = task.Status
			child.Message = task.Message

			var parent *taskNode
			if task.ParentID != nil {
				parent = getNode(*task.ParentID)
				child.Parent = parent
				parent.Children = append(parent.Children, child)
			}

			if task.ID == jobID {
				job = task
			}
		}

		depthOf := func(node *taskNode) int {
			d := 0
			for node != nil {
				node = node.Parent
				d++
			}
			return d
		}

		var rec func(idx int, node *taskNode)
		rec = func(idx int, node *taskNode) {

			prefix := ""
			if node.Parent != nil {
				if len(node.Children) > 0 || node.Status != taskapi.StatusRunning {
					prefix += "  "
				} else {
					// XXX const spinner = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
					const spinner = `/-\\|/-`
					frame := spinner[iteration%len(spinner)]
					prefix += string(frame) + " "
				}

				depth := depthOf(node)
				prefix += strings.Repeat("   ", depth-2)
				last := idx == len(node.Parent.Children)-1
				if last {
					prefix += "└─ "
				} else {
					prefix += "├─ "
				}
			}

			printf("%s%s %s %s", prefix, node.Name, node.Message, node.Status)
			for i, child := range node.Children {
				rec(i, child)
			}
		}
		root := getNode(job.ID)
		rec(0, root)

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
		iteration++
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
