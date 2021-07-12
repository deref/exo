package logrot

import (
	"context"

	"github.com/deref/exo/atom"
)

type service struct {
	statePath string
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

func (svc *service) AddLog(ctx context.Context, input *AddLogInput) (*AddLogOutput, error) {
	_, err := svc.swapState(func(state *State) error {
		state.Logs[input.ID] = LogState{}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// XXX kick off a worker, if one doesn't already exist.
	return &AddLogOutput{}, nil
}

func (svc *service) RemoveLog(ctx context.Context, input *RemoveLogInput) (*RemoveLogOutput, error) {
	_, err := svc.swapState(func(state *State) error {
		delete(state.Logs, input.ID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	// XXX kill the worker, if it exists.
	return &RemoveLogOutput{}, nil
}

func (svc *service) DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error) {
	state, err := svc.derefState()
	if err != nil {
		return nil, err
	}
	var output DescribeLogsOutput
	for id, description := range state.Logs {
		output.Logs = append(output.Logs, LogDescription{
			ID:          id,
			SourcePath:  description.SourcePath,
			LastEventAt: nil, // XXX set me to last line of file.
		})
	}
	return &output, nil
}

func (svc *service) GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error) {
	panic("TODO: GetEvents")
}
