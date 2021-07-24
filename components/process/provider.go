package process

type Provider struct {
	WorkspaceDir string
	VarDir       string
}

type Spec struct {
	Directory   string            `json:"directory"`
	Program     string            `json:"program"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

type State struct {
	Pid int `json:"pid"`
	// TODO: Store resolved program path & full effective environment.
	// Program string `json:"program"`
	// Environment map[string]string `json:"environment"`
}
