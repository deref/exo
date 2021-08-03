package process

import "github.com/deref/exo/providers/core"

type Process struct {
	core.Component
	Spec
	State

	SyslogPort int
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
	// TODO: Store resolved program path & full effective environment.
	// Program string `json:"program"`
	FullEnvironment map[string]string `json:"fullEnvironment"`
}
