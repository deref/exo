package server

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/deref/exo/logcol/api"
)

func validLogName(s string) bool {
	return s != "" // TODO: More validation?
}

func (lc *logCollector) AddLog(ctx context.Context, input *api.AddLogInput) (*api.AddLogOutput, error) {
	if !validLogName(input.Name) {
		return nil, fmt.Errorf("invalid log name: %q", input.Name)
	}
	if input.SourcePath == "" {
		return nil, errors.New("log source path is required")
	}
	state, err := lc.swapState(func(state *State) error {
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
	lc.startWorker(ctx, input.Name, state.Logs[input.Name])
	return &api.AddLogOutput{}, nil
}

func (lc *logCollector) RemoveLog(ctx context.Context, input *api.RemoveLogInput) (*api.RemoveLogOutput, error) {
	_, err := lc.swapState(func(state *State) error {
		if state.Logs != nil {
			delete(state.Logs, input.Name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	lc.stopWorker(input.Name)
	return &api.RemoveLogOutput{}, nil
}

func (lc *logCollector) DescribeLogs(context.Context, *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
	state, err := lc.derefState()
	if err != nil {
		return nil, err
	}
	var output api.DescribeLogsOutput
	output.Logs = []api.LogDescription{}
	for name, description := range state.Logs {
		output.Logs = append(output.Logs, api.LogDescription{
			Name:        name,
			SourcePath:  description.SourcePath,
			LastEventAt: nil, // XXX set me based on the last line of file.
		})
	}
	return &output, nil
}

func (lc *logCollector) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	// TODO: Handle Before & After pagination parameters.
	// TODO: Limit number of returned events.
	var output api.GetEventsOutput
	output.Events = []api.Event{}
	switch len(input.LogNames) {
	case 0:
		// nop.
	case 1:
		logName := input.LogNames[0]
		state, err := lc.derefState()
		if err != nil {
			return nil, err
		}
		logState, exists := state.Logs[logName]
		if exists {
			chunkIndex := 0 // TODO: Handle log rotation.
			f, err := os.Open(makeChunkPath(logState.SourcePath, chunkIndex))
			if err != nil {
				return nil, fmt.Errorf("opening %s source: %w", logName, err)
			}
			r := bufio.NewReader(f)
			line, isPrefix, err := r.ReadLine()
			if err != nil {
				return nil, fmt.Errorf("reading: %w", err)
			}
			// TODO: Do something better with lines that are too long.
			for isPrefix {
				// Skip remainder of line.
				line = append([]byte{}, line...)
				_, isPrefix, err = r.ReadLine()
				if err != nil {
					return nil, fmt.Errorf("reading: %w", err)
				}
			}
			fields := bytes.SplitN(line, []byte(" "), 3)
			if len(fields) != 3 {
				return nil, fmt.Errorf("invalid log line")
			}
			event := api.Event{
				LogName:   logName,
				SID:       string(fields[0]),
				Timestamp: string(fields[1]),
				Message:   string(fields[2]),
			}
			output.Events = append(output.Events, event)
		}
	default:
		return nil, fmt.Errorf("TODO: merge log streams")
	}
	return &output, nil
}

func (lc *logCollector) Collect(ctx context.Context, input *api.CollectInput) (*api.CollectOutput, error) {
	state, err := lc.derefState()
	if err != nil {
		return nil, err
	}
	for name, logState := range state.Logs {
		lc.startWorker(ctx, name, logState)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	return &api.CollectOutput{}, nil
}
