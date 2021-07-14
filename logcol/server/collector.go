package server

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strings"

	"github.com/deref/exo/logcol/api"
)

func validLogName(s string) bool {
	return s != "" // TODO: More validation?
}

func (lc *logCollector) AddLog(ctx context.Context, input *api.AddLogInput) (*api.AddLogOutput, error) {
	if !validLogName(input.Name) {
		return nil, fmt.Errorf("invalid log name: %q", input.Name)
	}
	if input.Source == "" {
		return nil, errors.New("log source path is required")
	}
	state, err := lc.swapState(func(state *State) error {
		if state.Logs == nil {
			state.Logs = make(map[string]LogState)
		}
		state.Logs[input.Name] = LogState{
			Source: input.Source,
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
			Source:      description.Source,
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
	for _, logName := range input.Logs {
		state, err := lc.derefState()
		if err != nil {
			return nil, err
		}
		logState, exists := state.Logs[logName]
		if exists {
			chunkIndex := 0 // TODO: Handle log rotation.
			f, err := os.Open(makeChunkPath(logState.Source, chunkIndex))
			if err != nil {
				return nil, fmt.Errorf("opening %s source: %w", logName, err)
			}
			r := bufio.NewReader(f)
			for {
				line, isPrefix, err := r.ReadLine()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, fmt.Errorf("reading: %w", err)
				}
				// TODO: Do something better with lines that are too long.
				for isPrefix {
					// Skip remainder of line.
					line = append([]byte{}, line...)
					_, isPrefix, err = r.ReadLine()
					if err == io.EOF {
						break
					}
					if err != nil {
						return nil, fmt.Errorf("reading: %w", err)
					}
				}
				fields := bytes.SplitN(line, []byte(" "), 3)
				if len(fields) != 3 {
					return nil, fmt.Errorf("invalid log line")
				}
				event := api.Event{
					Log:       logName,
					Sid:       string(fields[0]),
					Timestamp: string(fields[1]),
					Message:   string(fields[2]),
				}
				output.Events = append(output.Events, event)
			}
		}
	}
	sort.Sort(&eventsSorter{output.Events})
	return &output, nil
}

type eventsSorter struct {
	events []api.Event
}

func (iface *eventsSorter) Len() int {
	return len(iface.events)
}

func (iface *eventsSorter) Less(i, j int) bool {
	// TODO: Account for sequence ids.
	return strings.Compare(iface.events[i].Timestamp, iface.events[j].Timestamp) < 0
}

func (iface *eventsSorter) Swap(i, j int) {
	tmp := iface.events[i]
	iface.events[i] = iface.events[j]
	iface.events[j] = tmp
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
