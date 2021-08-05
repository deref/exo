package process

import "github.com/deref/exo/internal/providers/core"

type Process struct {
	core.Component
	Spec
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
	SupervisorPid int `json:"supervisorPid"`
	Pid           int `json:"pid"`
	// Program string `json:"program"`
	FullEnvironment map[string]string `json:"fullEnvironment"`
}
