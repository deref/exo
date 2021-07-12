package logrot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func NewHandler() http.Handler {
	statePath := "./var/logrot" // TODO: Configuration.
	return NewMux("/", NewService(statePath))
}

func validLogName(s string) bool {
	return s != "" // TODO: More validation?
}

func (svc *service) AddLog(ctx context.Context, input *AddLogInput) (*AddLogOutput, error) {
	if !validLogName(input.Name) {
		return nil, fmt.Errorf("invalid log name: %q", input.Name)
	}
	if input.SourcePath == "" {
		return nil, errors.New("log source path is required")
	}
	state, err := svc.swapState(func(state *State) error {
		if state.Logs == nil {
			state.Logs = make(map[string]LogState)
		}
		state.Logs[input.Name] = LogState{
			SourcePath: input.SourcePath,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	svc.startWorker(input.Name, state.Logs[input.Name])
	return &AddLogOutput{}, nil
}

func (svc *service) RemoveLog(ctx context.Context, input *RemoveLogInput) (*RemoveLogOutput, error) {
	_, err := svc.swapState(func(state *State) error {
		if state.Logs != nil {
			delete(state.Logs, input.Name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	svc.stopWorker(input.Name)
	return &RemoveLogOutput{}, nil
}

func (svc *service) DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error) {
	state, err := svc.derefState()
	if err != nil {
		return nil, err
	}
	var output DescribeLogsOutput
	output.Logs = []LogDescription{}
	for name, description := range state.Logs {
		output.Logs = append(output.Logs, LogDescription{
			Name:        name,
			SourcePath:  description.SourcePath,
			LastEventAt: nil, // XXX set me based on the last line of file.
		})
	}
	return &output, nil
}

func (svc *service) GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error) {
	panic("TODO: GetEvents")
}

func (svc *service) CollectLogs(context.Context, *CollectLogsInput) (*CollectLogsOutput, error) {
	state, err := svc.derefState()
	if err != nil {
		return nil, err
	}
	for name, logState := range state.Logs {
		svc.startWorker(name, logState)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	return &CollectLogsOutput{}, nil
}
