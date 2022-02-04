package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/deref/exo/internal/api"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/util/term"
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
			err = api.Query(ctx, svc, &q, map[string]interface{}{
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
	ID       string
	ParentID *string
	JobID    string
	Finished *string
	Created  string
	Status   string
	Label    string
	Message  string
	Progress *progressFragment
}

type taskNode struct {
	Parent   *taskNode
	ID       string
	Created  string
	Status   string
	Label    string
	Message  string
	Progress *progressFragment
	Children []*taskNode
}

type progressFragment struct {
	Percent float64
}

type jobPrinter struct {
	Spinner   []string
	Iteration int
	ShowJobID bool
}

func (jp *jobPrinter) printTree(w io.Writer, tasks []taskFragment) {
	// TODO: watchJob calls printTree in a loop, so each go around calls
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
		child.Label = task.Label
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

	// Measure labels.
	maxPrefixW := 8
	for _, task := range tasks {
		node := getNode(task.ID)
		depth := depthOf(node)
		statusW := 2
		indentW := 3
		padRight := 2
		width := statusW + (depth-1)*indentW + term.VisualLength(node.Label) + padRight
		if jp.ShowJobID {
			width += 1 + term.VisualLength(node.ID)
		}
		if maxPrefixW < width {
			maxPrefixW = width
		}
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

		label := node.Label
		if jp.ShowJobID {
			label += " " + node.ID
		}
		prefix += label

		// Padding between prefix and message to align message column.
		prefixW := term.VisualLength(prefix)
		alignMessage := ""
		padPrefix := maxPrefixW - prefixW
		if padPrefix > 0 {
			alignMessage = strings.Repeat(" ", padPrefix)
		}

		suffix := ""
		if node.Progress != nil {
			percent := node.Progress.Percent
			if percent < 100 {
				suffix = fmt.Sprintf("  %2d %% ", int(percent))
			}
		}

		message := node.Message

		// Truncate message.
		maxMessageW := 50
		if termW > 0 {
			suffixW := term.VisualLength(suffix)
			maxMessageW = termW - maxPrefixW - suffixW
		}
		message = term.TrimToVisualLength(message, maxMessageW)

		// Right align suffix.
		messageW := term.VisualLength(message)
		alignSuffix := strings.Repeat(" ", maxMessageW-messageW)

		fmt.Fprintf(w, "%s%s%s%s%s\n", prefix, alignMessage, message, alignSuffix, suffix)
		for i, child := range node.Children {
			rec(i, child)
		}
	}
	for _, job := range jobNodes {
		rec(0, job)
	}
}
