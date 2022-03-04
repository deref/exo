package daemon

import "github.com/deref/exo/internal/providers/sdk"

type ComponentConfig struct {
	sdk.ComponentConfig
	Model
}

type Model struct {
	Spec
	State
}

type Spec struct {
	Program     string            `json:"program"`
	Arguments   []string          `json:"arguments,omitempty"`
	Directory   string            `json:"directory,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

type State struct {
	// TODO: Transition
}

type Controller struct{}
