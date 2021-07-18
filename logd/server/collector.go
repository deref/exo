package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/deref/exo/atom"
	"github.com/deref/exo/gensym"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/logd/store"
	"github.com/deref/exo/util/mathutil"
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
		idGen:  gensym.NewULIDGenerator(ctx),
		debug:  cfg.Debug,
	}
}

type LogCollector struct {
	varDir string
	state  atom.Atom
	idGen  *gensym.ULIDGenerator
	debug  bool
	store  store.Store
	wg     sync.WaitGroup

	mx      sync.Mutex
	workers map[string]*collectorWorker
}

type collectorWorker struct {
	Worker
	stop func()
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

func (lc *LogCollector) ensureWorker(ctx context.Context, logName string, state LogState) {
	lc.mx.Lock()
	defer lc.mx.Unlock()
	if lc.workers == nil {
		// No worker support in peer mode.
		return
	}
	lc.startWorker(ctx, logName, state)
}

func (lc *LogCollector) startWorker(ctx context.Context, logName string, state LogState) {
	wkr, exists := lc.workers[logName]
	if exists {
		return
	}
	ctx, stop := context.WithCancel(ctx)
	wkr = &collectorWorker{
		Worker: Worker{
			Source: state.Source,
			Sink:   lc.store.GetLog(logName),
			Debug:  lc.debug,
		},
		stop: stop,
	}
	lc.workers[logName] = wkr

	done := make(chan struct{})
	lc.wg.Add(1)

	go func() {
		defer wkr.debugf("run done")
		defer lc.wg.Done()
		err := wkr.Run(ctx)
		if err != nil {
			// TODO: Panic instead.
			fmt.Fprintf(os.Stderr, "worker run error: %v\n", err)
		}
		close(done)
	}()

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

func (lc *LogCollector) stopWorker(logName string) {
	lc.mx.Lock()
	defer lc.mx.Unlock()

	wkr := lc.workers[logName]
	if wkr == nil {
		return
	}

	wkr.stop()
}

func (lc *LogCollector) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
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
			LastEventAt: lc.store.GetLog(name).GetLastEventAt(ctx),
		})
	}
	return &output, nil
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
	for _, logName := range logs {
		log := lc.store.GetLog(logName)
		logEvents, err := log.GetEvents(ctx, input.After, limit)
		if err != nil {
			return nil, fmt.Errorf("getting %q events: %w", logName, err)
		}
		events = append(events, logEvents...)
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
