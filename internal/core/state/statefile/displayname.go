package statefile

import (
	"fmt"
	"path/filepath"
	"strings"
)

type displayNameBuilder struct {
	tree      reversePathTree
	separator string
}

func newDisplayNameBuilder(separator string) *displayNameBuilder {
	return &displayNameBuilder{
		tree:      make(reversePathTree),
		separator: separator,
	}
}

type reversePathTree map[string]reversePathTree

func (b *displayNameBuilder) dump() {
	b.tree.dump("")
}

func (tree reversePathTree) dump(prefix string) {
	for name, child := range tree {
		fmt.Printf("%s%q\n", prefix, name)
		child.dump(prefix + "  ")
	}
}

func (b *displayNameBuilder) pathParts(path string) []string {
	return strings.Split(filepath.Clean(path), b.separator)
}

func (b *displayNameBuilder) AddPath(path string) {
	parts := b.pathParts(path)
	tree := b.tree
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		child := tree[part]
		if child == nil {
			child = make(reversePathTree)
			tree[part] = child
		}
		tree = child
	}
}

func (b *displayNameBuilder) GetDisplayName(path string) string {
	parts := b.pathParts(path)
	tree := b.tree
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		child := tree[part]
		if len(child) <= 1 {
			return strings.Join(parts[i:], b.separator)
		}
	}
	// Fallback, but shouldn't be reachable if builder is properly initialized.
	return path
}
