package logcol

import (
	"sync"

	"github.com/deref/exo/atom"
)

type State struct {
	Logs map[string]LogState `json:"logs"`
}

type LogState struct {
	SourcePath string `json:"sourcePath"`
}

func NewService() Service {
	statePath := "./var/logcol" // TODO: Configuration.
	return &service{
		state:   atom.NewFileAtom(statePath, atom.CodecJSON),
		workers: make(map[string]*worker),
	}
}

type service struct {
	state atom.Atom

	mx      sync.Mutex
	workers map[string]*worker
}

func (svc *service) derefState() (*State, error) {
	var state State
	err := svc.state.Deref(&state)
	return &state, err
}

func (svc *service) swapState(f func(state *State) error) (*State, error) {
	var state State
	err := svc.state.Swap(&state, func() error {
		return f(&state)
	})
	return &state, err
}
