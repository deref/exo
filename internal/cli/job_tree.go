package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/term"
	"github.com/deref/rgbterm"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobTreeCmd)
}

// TODO: Scope flag: "all", "stack", etc. Default to current stack.

var jobTreeCmd = &cobra.Command{
	Use:   "tree [job-id...]",
	Short: "Lists jobs with a tree of their tasks",
	Long: `Lists jobs with a tree of their tasks.

If no job ids are provided, lists all jobs in the scope.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var tasks []taskFragment
		var err error
		// TODO: Filter by scope.
		if len(args) == 0 {
			var q struct {
				Tasks []taskFragment `graphql:"allTasks"`
			}
			err = api.Query(ctx, svc, &q, nil)
			tasks = q.Tasks
		} else {
			var q struct {
				Tasks []taskFragment `graphql:"tasksByJobIds(jobIds: $jobIds)"`
			}
			err = api.Query(ctx, svc, &q, map[string]any{
				"jobIds": args,
			})
			tasks = q.Tasks
		}
		if err != nil {
			return fmt.Errorf("querying tasks: %w", err)
		}

		jp := &jobPrinter{
			ShowJobID: len(args) != 1,
		}
		jp.printTree(os.Stdout, tasks)
		return nil
	},
}

type taskFragment struct {
	ID         string
	ParentID   *string
	JobID      string
	Finished   *string
	Created    string
	Label      string
	Message    string
	Started    *scalars.Instant
	Completed  *scalars.Instant
	Successful *bool
	Error      *string
	Progress   *progressFragment
}

const (
	taskStatusQueue = "queue"
	taskStatusRun   = "run"
	taskStatusWait  = "wait"
	taskStatusOK    = "ok"
	taskStatusErr   = "err"
)

func classifyTaskStatus(task taskFragment) string {
	if task.Error != nil {
		return taskStatusErr
	}
	if task.Started == nil {
		return taskStatusQueue
	}
	if task.Finished == nil {
		return taskStatusRun
	}
	if task.Completed == nil {
		return taskStatusWait
	}
	return taskStatusOK
}

type taskNode struct {
	Parent     *taskNode
	ID         string
	Created    string
	Status     string
	Label      string
	Message    string
	Progress   *progressFragment
	ChildIndex int
	Children   []*taskNode
}

type progressFragment struct {
	Percent float64
}

type jobPrinter struct {
	Spinner            []string
	Iteration          int
	ShowJobID          bool
	CollapseSuccessful bool
}

func jobTreeString(jp *jobPrinter, tasks []taskFragment) string {
	var sb strings.Builder
	jp.printTree(&sb, tasks)
	return sb.String()
}

func (jp *jobPrinter) printTree(w io.Writer, tasks []taskFragment) {
	builder := term.NewTreeBuilder()
	for _, task := range tasks {
		node := &term.TreeNode{
			ID: task.ID,
		}
		if task.ParentID != nil {
			node.ParentID = *task.ParentID
		}

		status := classifyTaskStatus(task)
		switch status {
		case taskStatusQueue:
			node.Label = rgbterm.FgString("·", 0, 123, 211)
		case taskStatusRun:
			if len(jp.Spinner) > 0 {
				offset := jp.Iteration + node.ChildIndex
				frame := jp.Spinner[offset%len(jp.Spinner)]
				node.Label = rgbterm.FgString(frame, 172, 66, 199)
			} else {
				node.Label = " "
			}
		case taskStatusWait:
			node.Label = rgbterm.FgString("○", 0, 123, 211)
		case taskStatusOK:
			node.Label = rgbterm.FgString("✓", 28, 196, 22)
		case taskStatusErr:
			node.Label = rgbterm.FgString("⨯", 215, 55, 30)
		default:
			panic(fmt.Sprintf("unexpected task status: %q", status))
		}

		node.Label += " " + task.Label
		if jp.ShowJobID {
			node.Label += " " + node.ID
		}

		if task.Progress != nil {
			percent := task.Progress.Percent
			if percent < 100 {
				node.Suffix = fmt.Sprintf("  %2d %% ", int(percent))
			}
		}

		if jp.CollapseSuccessful && task.Successful != nil && *task.Successful == true {
			node.HideChildren = true
		}

		builder.AddNode(node)
	}

	trees := builder.Build()
	for _, tree := range trees {
		term.PrintTree(w, tree)
	}
}
