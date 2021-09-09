package process

import "github.com/deref/exo/internal/providers/core"

type Process struct {
	core.ComponentBase
	Spec
	State

	SyslogPort uint
}

type Spec struct {
	Directory                  string            `json:"directory,omitempty"`
	Program                    string            `json:"program"`
	Arguments                  []string          `json:"arguments,omitempty"`
	Environment                map[string]string `json:"environment"`
	ShutdownGracePeriodSeconds *int              `json:"shutdownGracePeriodSeconds,omitempty"`
}

type State struct {
	Pgid          int `json:"pgid"`
	SupervisorPid int `json:"supervisorPid"`
	Pid           int `json:"pid"`
	// Program string `json:"program"`
	FullEnvironment map[string]string `json:"fullEnvironment"`
}

func (state *State) clear() {
	*state = State{}
}
