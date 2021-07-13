package server

import (
	"sync"

	"github.com/deref/exo/atom"
	"github.com/deref/exo/logcol/api"
)

type State struct {
	Logs map[string]LogState `json:"logs"`
}

type LogState struct {
	Source string `json:"source"`
}

func NewLogCollector() api.LogCollector {
	statePath := "./var/logcol.json" // TODO: Configuration.
	return &logCollector{
		state:   atom.NewFileAtom(statePath, atom.CodecJSON),
		workers: make(map[string]*worker),
	}
}

type logCollector struct {
	state atom.Atom

	mx      sync.Mutex
	workers map[string]*worker
}

func (lc *logCollector) derefState() (*State, error) {
	var state State
	err := lc.state.Deref(&state)
	return &state, err
}

func (lc *logCollector) swapState(f func(state *State) error) (*State, error) {
	var state State
	err := lc.state.Swap(&state, func() error {
		return f(&state)
	})
	return &state, err
}
