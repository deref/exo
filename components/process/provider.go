package process

type Provider struct {
	ProjectDir string
	VarDir     string
}

type spec struct {
	Directory   string            `json:"directory"`
	Command     string            `json:"command"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

type state struct {
	Pid int `json:"pid"`
	// TODO: Store resolved command path & full effective environment.
	// Command string `json:"command"`
	// Environment string `json:"environment"`
}
