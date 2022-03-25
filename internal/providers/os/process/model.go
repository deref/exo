package process

import "github.com/deref/exo/internal/providers/sdk"

type Component struct {
	sdk.ComponentConfig
	Model
}

type Model struct {
	Spec
	State
}

type Spec struct {
	Program     string            `json:"program,omitempty"`
	Arguments   []string          `json:"arguments,omitempty"`
	Directory   string            `json:"directory,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

type State struct {
	ProgramPath string `json:"programPath,omitempty"`
	Pid         *int   `json:"pid,omitempty"`
}
