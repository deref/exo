package process

type Provider struct {
	ProjectDir string
	VarDir     string
}

type spec struct {
	Directory string   `json:"directory"`
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

type state struct {
	Pid int `json:"pid"`
}
