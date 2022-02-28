package term

import (
	"fmt"
	"io"
	"strings"
)

type TreeBuilder struct {
	byID map[string]*TreeNode
	all  []*TreeNode
}

func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		byID: make(map[string]*TreeNode),
	}
}

func (b *TreeBuilder) AddNode(node *TreeNode) {
	b.byID[node.ID] = node
	b.all = append(b.all, node)
}

func (b *TreeBuilder) Build() (roots []*TreeNode) {
	for _, node := range b.all {
		if node.ParentID == "" {
			roots = append(roots, node)
		} else {
			parent := b.byID[node.ParentID]
			node.Parent = parent
			node.ChildIndex = len(parent.Children)
			parent.Children = append(parent.Children, node)
		}
	}

	var visit func(depth int, node *TreeNode)
	visit = func(depth int, node *TreeNode) {
		indentWidth := 3
		padRight := 2
		width := depth*indentWidth + VisualLength(node.Label) + padRight
		for _, child := range node.Children {
			visit(depth+1, child)
			childWidth := child.LabelWidth
			if childWidth > width {
				width = childWidth
			}
		}
		node.LabelWidth = width
	}
	for _, root := range roots {
		visit(0, root)
	}

	return
}

type TreeNode struct {
	ID           string
	ParentID     string
	Label        string // Left-most column, part of the tree shape.
	Content      string // Second column, aligned after the labels.
	Suffix       string // Right-aligned, third column.
	HideChildren bool

	// Fields set by TreeBuilder.
	Parent     *TreeNode
	ChildIndex int
	Children   []*TreeNode
	LabelWidth int // Max visual width of own and children's label with padding.
}

func PrintTree(w io.Writer, node *TreeNode) {
	printTree(w, node, node)
}

func printTree(w io.Writer, root *TreeNode, node *TreeNode) {
	// TODO: watchJob calls printTree in a loop, so each go around calls
	// term.GetSize(), it would be more efficient to listen to terminal size
	// change events.
	termW, _ := GetSize()

	var prefixBuilder strings.Builder
	printTracks(&prefixBuilder, node, "└─ ", "├─ ")
	prefixBuilder.WriteString(node.Label)
	prefix := prefixBuilder.String()

	prefixWidth := VisualLength(prefix)
	alignContent := ""
	padPrefix := root.LabelWidth - prefixWidth
	if padPrefix > 0 {
		alignContent = strings.Repeat(" ", padPrefix)
	}

	suffix := node.Suffix

	// Truncate content.
	content := node.Content
	maxContentW := 50
	if termW > 0 {
		suffixW := VisualLength(suffix)
		maxContentW = termW - root.LabelWidth - suffixW
	}
	content = TrimToVisualLength(content, maxContentW)

	// Right align suffix.
	contentW := VisualLength(content)
	alignSuffix := strings.Repeat(" ", maxContentW-contentW)

	fmt.Fprintf(w, "%s%s%s%s%s\n", prefix, alignContent, content, alignSuffix, suffix)
	if !node.HideChildren {
		for _, child := range node.Children {
			printTree(w, root, child)
		}
	}
}

func printTracks(w io.Writer, node *TreeNode, end, more string) {
	if node.Parent == nil {
		return
	}
	last := node.ChildIndex == len(node.Parent.Children)-1
	printTracks(w, node.Parent, "   ", "│  ")
	if last {
		io.WriteString(w, end)
	} else {
		io.WriteString(w, more)
	}
}
