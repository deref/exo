package logrot

type State struct {
	Logs map[string]LogState `json:"logs"`
}

type LogState struct {
	SourcePath string `json:"sourcePath"`
}
