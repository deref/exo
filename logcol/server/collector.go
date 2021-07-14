package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deref/exo/atom"
	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logcol/api"
	badger "github.com/dgraph-io/badger/v3"
)

func NewLogCollector() *LogCollector {
	varDir := "./var" // TODO: Configuration?
	statePath := filepath.Join(varDir, "logcol.json")
	return &LogCollector{
		varDir: varDir,
		state:  atom.NewFileAtom(statePath, atom.CodecJSON),
	}
}

type LogCollector struct {
	varDir string
	state  atom.Atom
	db     *badger.DB

	mx      sync.Mutex
	workers map[string]*worker
}

func validLogName(s string) bool {
	return s != "" // TODO: More validation. Cannot have internal null bytes.
}

func (lc *LogCollector) AddLog(ctx context.Context, input *api.AddLogInput) (*api.AddLogOutput, error) {
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
	lc.ensureWorker(ctx, input.Name, state.Logs[input.Name])
	return &api.AddLogOutput{}, nil
}

func (lc *LogCollector) RemoveLog(ctx context.Context, input *api.RemoveLogInput) (*api.RemoveLogOutput, error) {
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

func (lc *LogCollector) DescribeLogs(context.Context, *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
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
			LastEventAt: lc.getLastEventAt(name),
		})
	}
	return &output, nil
}

func (lc *LogCollector) getLastEventAt(name string) *string {
	fmt.Println("get last event", name)
	var timestamp uint64
	if err := lc.db.View(func(txn *badger.Txn) error {
		prefix := append([]byte(name), 0)
		it := txn.NewIterator(badger.IteratorOptions{
			Prefix:  prefix,
			Reverse: true,
		})
		defer it.Close()
		if !it.Valid() {
			fmt.Println("not valid")
			return nil
		}
		fmt.Print("HERE")
		return it.Item().Value(func(bs []byte) error {
			timestamp = binary.BigEndian.Uint64(bs)
			return nil
		})
	}); err != nil {
		fmt.Println("view err:", err)
		return nil
	}
	if timestamp != 0 {
		s := chrono.NanoToIso(int64(timestamp))
		return &s
	}
	return nil
}

func (lc *LogCollector) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	return nil, nil
	/*
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
				r := bufio.NewReaderSize(f, api.MaxEventSize)
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
					event, err := parseEvent(logName, line)
					if err != nil {
						return nil, err
					}
					output.Events = append(output.Events, *event)
				}
			}
		}
		sort.Sort(&eventsSorter{output.Events})
		return &output, nil
	*/
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
