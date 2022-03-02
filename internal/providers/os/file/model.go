package file

import "github.com/deref/exo/internal/scalars"

type Model struct {
	Spec
	State
}

type Spec struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type State struct {
	Size     *string          `json:"size,omitempty"` // String for int64 support.
	Modified *scalars.Instant `json:"modified,omitempty"`
}

type Controller struct{}
