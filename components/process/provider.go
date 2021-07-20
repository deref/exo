package process

type Provider struct {
	WorkspaceDir string
	VarDir       string
	Fifofum      FifofumConfig
}

type FifofumConfig struct {
	Path string
	Args []string
}

var FifofumDevConfig = FifofumConfig{
	Path: "go",
	Args: []string{"run", "./cmd/fifofum"},
}

type spec struct {
	Directory   string            `json:"directory"`
	Program     string            `json:"program"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

type state struct {
	Pid int `json:"pid"`
	// TODO: Store resolved program path & full effective environment.
	// Program string `json:"program"`
	// Environment string `json:"environment"`
}
