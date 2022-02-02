package file

type Model struct {
	Spec
	State
}

type Spec struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type State struct {
	// TODO: mtime, etc.
}

type Controller struct{}
