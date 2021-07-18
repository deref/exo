package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/deref/exo/atom"
	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/util/mathutil"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/oklog/ulid/v2"
)

const (
	defaultEventsPerRequest = 500
)

type Config struct {
	VarDir string
	Debug  bool
}

func NewLogCollector(ctx context.Context, cfg *Config) *LogCollector {
	statePath := filepath.Join(cfg.VarDir, "logd.json")
	return &LogCollector{
		varDir: cfg.VarDir,
		state:  atom.NewFileAtom(statePath, atom.CodecJSON),
		idGen:  newIdGen(ctx),
		debug:  cfg.Debug,
	}
}

type LogCollector struct {
	varDir string
	state  atom.Atom
	idGen  *idGen
	debug  bool
	db     *badger.DB
	wg     sync.WaitGroup

	mx      sync.Mutex
	workers map[string]*worker
}

func (lc *LogCollector) debugf(format string, v ...interface{}) {
	if lc.debug {
		fmt.Fprintln(os.Stderr, "collector", fmt.Errorf(format, v...))
	}
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
				if err := validateVersion(val[versionOffset]); err != nil {
					return err
				}
				timestamp = binary.BigEndian.Uint64(val[timestampOffset : timestampOffset+timestampLen])
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
				id, err := ulid.Parse(strings.ToUpper(input.After))
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
		end := mathutil.IntMin(limit, len(events))
		events = events[:end]

		cursor = events[len(events)-1].ID
	}

	return &api.GetEventsOutput{
		Events: events,
		Cursor: cursor,
	}, nil
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
