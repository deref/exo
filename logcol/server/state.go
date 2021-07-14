package server

type State struct {
	Logs map[string]LogState `json:"logs"`
}

type LogState struct {
	Source string `json:"source"`
}

func (lc *LogCollector) derefState() (*State, error) {
	var state State
	err := lc.state.Deref(&state)
	return &state, err
}

func (lc *LogCollector) swapState(f func(state *State) error) (*State, error) {
	var state State
	err := lc.state.Swap(&state, func() error {
		return f(&state)
	})
	return &state, err
}
