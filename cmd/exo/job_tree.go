package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/deref/exo/internal/core/api"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobTreeCmd)
}

var jobTreeCmd = &cobra.Command{
	Use:   "tree [job-ids...]",
	Short: "Lists jobs with a tree of their tasks",
	Long: `Lists jobs with a tree of their tasks.

If job-ids are provided, lists all jobs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		kernel := cl.Kernel()

		input := &api.DescribeTasksInput{}
		if len(args) > 0 {
			input.JobIDs = args
		}
		output, err := kernel.DescribeTasks(ctx, input)
		if err != nil {
			return fmt.Errorf("describing tasks: %w", err)
		}

		jp := &jobPrinter{
			ShowJobID: len(args) != 1,
		}
		jp.printTree(os.Stdout, output.Tasks)
		return nil
	},
}

type taskNode struct {
	Parent   *taskNode
	ID       string
	Created  string
	Status   string
	Name     string
	Message  string
	Progress *api.TaskProgress
	Children []*taskNode
}

type jobPrinter struct {
	Spinner   []string
	Iteration int
	ShowJobID bool
}

func (jp *jobPrinter) printTree(w io.Writer, tasks []api.TaskDescription) {
	// TODO: watchJobs calls printTree in a loop, so each go around calls
	// term.GetSize(), it would be more efficient to listen to terminal size
	// change events.
	termW, _ := term.GetSize()

	nodes := make(map[string]*taskNode)
	getNode := func(id string) *taskNode {
		node := nodes[id]
		if node == nil {
			node = &taskNode{}
			nodes[id] = node
		}
		return node
	}

	var jobNodes []*taskNode
	for _, task := range tasks {
		child := getNode(task.ID)
		child.ID = task.ID
		child.Created = task.Created
		child.Name = task.Name
		child.Status = task.Status
		child.Message = task.Message
		child.Progress = task.Progress

		var parent *taskNode
		if task.ParentID == nil {
			jobNodes = append(jobNodes, child)
		} else {
			parent = getNode(*task.ParentID)
			child.Parent = parent
			parent.Children = append(parent.Children, child)
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
			switch node.Status {
			case taskapi.StatusPending:
				prefix += rgbterm.FgString("·", 0, 123, 211)
			case taskapi.StatusSuccess:
				prefix += rgbterm.FgString("✓", 28, 196, 22)
			case taskapi.StatusFailure:
				prefix += rgbterm.FgString("⨯", 215, 55, 30)
			case taskapi.StatusRunning:
				if len(jp.Spinner) > 0 {
					offset := jp.Iteration + idx
					frame := jp.Spinner[offset%len(jp.Spinner)]
					prefix += rgbterm.FgString(frame, 172, 66, 199)
				}
			}
			prefix += " "

			depth := depthOf(node)
			prefix += strings.Repeat("│  ", depth-2)
			last := idx == len(node.Parent.Children)-1
			if last {
				prefix += "└─ "
			} else {
				prefix += "├─ "
			}
		}

		label := node.Name
		if jp.ShowJobID {
			label += " " + node.ID
		}
		if node.Parent != nil {
			prefix += " "
		}
		prefix += label + " "

		suffix := ""
		if node.Progress != nil {
			progress := float64(node.Progress.Current) / float64(node.Progress.Total)
			if progress < 1 {
				suffix = fmt.Sprintf("  %2d %% ", int(progress*100.0))
			}
		}

		message := node.Message

		// Truncate message.
		maxMessageW := 50
		if termW > 0 {
			prefixW := term.VisualLength(prefix)
			suffixW := term.VisualLength(suffix)
			maxMessageW = termW - prefixW - suffixW
		}
		message = term.TrimToVisualLength(message, maxMessageW)

		// Right align suffix.
		messageW := term.VisualLength(message)
		spacer := strings.Repeat(" ", maxMessageW-messageW)

		fmt.Fprintf(w, "%s%s%s%s\n", prefix, message, spacer, suffix)
		for i, child := range node.Children {
			rec(i, child)
		}
	}
	for _, job := range jobNodes {
		rec(0, job)
	}
}
