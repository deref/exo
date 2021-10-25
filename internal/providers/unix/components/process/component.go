package process

import "github.com/deref/exo/internal/providers/core"

type Process struct {
	core.ComponentBase
	State

	SyslogPort uint
}

type Spec struct {
	Directory                  string            `json:"directory"`
	Program                    string            `json:"program"`
	Arguments                  []string          `json:"arguments"`
	Environment                map[string]string `json:"environment"`
	ShutdownGracePeriodSeconds *int              `json:"shutdownGracePeriodSeconds"`
}

type State struct {
	Directory                  string            `json:"directory"`
	Program                    string            `json:"program"`
	Arguments                  []string          `json:"arguments"`
	Environment                map[string]string `json:"environment"`
	ShutdownGracePeriodSeconds *int              `json:"shutdownGracePeriodSeconds"`

	Pgid            int               `json:"pgid"`
	SupervisorPid   int               `json:"supervisorPid"`
	Pid             int               `json:"pid"`
	FullEnvironment map[string]string `json:"fullEnvironment"`
}

func (state *State) reset() {
	state.Pgid = 0
	state.SupervisorPid = 0
	state.Pid = 0
	state.FullEnvironment = nil
}
