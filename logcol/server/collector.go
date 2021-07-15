package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"github.com/deref/exo/atom"
	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logcol/api"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/oklog/ulid/v2"
)

const (
	defaultEventsPerRequest = 500
)

type Config struct {
	VarDir string
}

func NewLogCollector(cfg *Config) *LogCollector {
	statePath := filepath.Join(cfg.VarDir, "logcol.json")
	return &LogCollector{
		varDir: cfg.VarDir,
		state:  atom.NewFileAtom(statePath, atom.CodecJSON),
		idGen:  newIdGen(),
	}
}

type LogCollector struct {
	varDir string
	state  atom.Atom
	idGen  *idGen
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
	prefix := append([]byte(name), 0)
	var timestamp uint64
	err := lc.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 1
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		it.Seek(append([]byte(name), 255))
		if it.ValidForPrefix(prefix) {
			item := it.Item()
			return item.Value(func(val []byte) error {
				timestamp = binary.BigEndian.Uint64(val[9 : 9+8])
				return nil
			})
		}
		return nil
	})
	if err != nil {
		return nil
	}

	if timestamp != 0 {
		s := chrono.NanoToIso(int64(timestamp))
		return &s
	}
	return nil
}

func (lc *LogCollector) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	logs := input.Logs
	if logs == nil {
		state, err := lc.derefState()
		if err != nil {
			return nil, fmt.Errorf("getting state: %w", err)
		}

		for name := range state.Logs {
			logs = append(logs, name)
		}
	}

	// TODO: Allow override via input.
	limit := defaultEventsPerRequest

	// TODO: Sorted merge.
	events := []api.Event{}
	for _, log := range logs {
		if err := lc.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			prefix := append([]byte(log), 0)
			start := prefix
			if input.After != "" {
				id, err := ulid.Parse(input.After)
				if err != nil {
					return fmt.Errorf("parsing cursor: %w", err)
				}
				idBytes, _ := id.MarshalBinary() // Cannot fail
				start = append(start, incrementBytes(idBytes)...)
			}

			eventsProcessed := 0
			for it.Seek(start); it.ValidForPrefix(prefix) && eventsProcessed < limit; it.Next() {
				item := it.Item()
				key := item.Key()
				if err := item.Value(func(val []byte) error {
					evt, err := eventFromEntry(log, key, val)
					if err != nil {
						return err
					}

					events = append(events, evt)
					return nil
				}); err != nil {
					return err
				}
				eventsProcessed++
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("scanning index: %w", err)
		}
	}

	var cursor string
	if len(events) > 0 {
		sort.Sort(&eventsSorter{events})
		events = events[:limit]

		cursor = events[len(events)-1].ID
	}

	return &api.GetEventsOutput{
		Events: events,
		Cursor: cursor,
	}, nil
}

func eventFromEntry(log string, key, val []byte) (api.Event, error) {
	// Parse key as (logName, null, id).
	logNameOffset := 0
	logNameLen := len(log)
	idOffset := logNameOffset + logNameLen + 1
	id, err := parseID(key[idOffset:])
	if err != nil {
		return api.Event{}, fmt.Errorf("parsing id: %w", err)
	}

	// Create value as (version, timestamp, message).
	// Version is used so that we can change the value format without rebuilding the database.
	versionOffset := 0
	versionLen := 1
	timestampOffset := versionOffset + versionLen
	timestampLen := 8
	messageOffset := timestampOffset + timestampLen

	version := val[versionOffset]
	if version != 1 {
		return api.Event{}, fmt.Errorf("unsupported event version: %d; database may have been written with a newer version of exo.", version)
	}

	tsNano := binary.BigEndian.Uint64(val[timestampOffset : timestampOffset+timestampLen])
	message := string(val[messageOffset:])

	return api.Event{
		ID:        id,
		Log:       log,
		Timestamp: chrono.NanoToIso(int64(tsNano)),
		Message:   message,
	}, nil
}

// incrementBytes returns a byte slice that is incremented by 1 bit.
// If `val` is not already only 255-valued bytes, then it is mutated and returned.
// Otherwise, a new slice is allocated and returned.
func incrementBytes(val []byte) []byte {
	for idx := len(val) - 1; idx >= 0; idx-- {
		byt := val[idx]
		if byt == 255 {
			val[idx] = 0
		} else {
			val[idx] = byt + 1
			return val
		}
	}

	// Still carrying from previously most significant byte, so add a new 1-valued byte.
	newVal := make([]byte, len(val)+1)
	newVal[0] = 1
	return newVal
}

type eventsSorter struct {
	events []api.Event
}

func (iface *eventsSorter) Len() int {
	return len(iface.events)
}

func (iface *eventsSorter) Less(i, j int) bool {
	iId := ulid.MustParse(iface.events[i].ID)
	jId := ulid.MustParse(iface.events[j].ID)
	return iId.Compare(jId) < 0
}

func (iface *eventsSorter) Swap(i, j int) {
	tmp := iface.events[i]
	iface.events[i] = iface.events[j]
	iface.events[j] = tmp
}
