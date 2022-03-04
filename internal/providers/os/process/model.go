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
	Program     string            `json:"program"`
	Arguments   []string          `json:"arguments"`
	Directory   string            `json:"directory"`
	Environment map[string]string `json:"environment"`
}

type State struct {
	ProgramPath string `json:"programPath"`
	Pid         *int   `json:"pid"`
}

type Controller struct {
	sdk.ResourceComponentController
}
