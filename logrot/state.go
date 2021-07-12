package logrot

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

func NewService(statePath string) Service {
	return &service{
		statePath: statePath,
		workers:   make(map[string]*worker),
	}
}

type service struct {
	mx        sync.Mutex
	statePath string
	workers   map[string]*worker
}

func (svc *service) derefState() (*State, error) {
	var state State
	err := atom.DerefJSON(svc.statePath, &state)
	return &state, err
}

func (svc *service) swapState(f func(state *State) error) (*State, error) {
	var state State
	err := atom.SwapJSON(svc.statePath, &state, func() error {
		return f(&state)
	})
	return &state, err
}
