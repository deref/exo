package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/deref/exo/gensym"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/logd/store"
	"github.com/deref/exo/logd/store/badger"
	"github.com/deref/exo/util/agent"
	"github.com/deref/exo/util/binaryutil"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/mathutil"
	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
)

const (
	defaultLimit = 500
)

type Config struct {
	VarDir string
	Debug  bool
}

func NewLogCollector(ctx context.Context, cfg *Config) *LogCollector {
	eg, ctx := errgroup.WithContext(ctx)
	return &LogCollector{
		varDir:    cfg.VarDir,
		idGen:     gensym.NewULIDGenerator(ctx),
		debug:     cfg.Debug,
		statePath: filepath.Join(cfg.VarDir, "logd.json"),
		eg:        eg,
		agent:     agent.NewAgent(300),
	}
}

type LogCollector struct {
	varDir string
	idGen  *gensym.ULIDGenerator
	debug  bool

	state
	statePath string

	store   store.Store
	eg      *errgroup.Group
	agent   *agent.Agent
	workers map[string]*collectorWorker
}

type state struct {
	Logs map[string]*logState `json:"logs"`
}

type logState struct {
	Source string `json:"source"`
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

func (lc *LogCollector) loadState() error {
	return jsonutil.UnmarshalFile(lc.statePath, &lc.state)
}

func (lc *LogCollector) saveState() error {
	return jsonutil.MarshalFile(lc.statePath, &lc.state)
}

func (lc *LogCollector) Run(ctx context.Context) error {
	lc.workers = make(map[string]*collectorWorker)

	if err := lc.loadState(); err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	logsDir := filepath.Join(lc.varDir, "logs")
	var err error
	lc.store, err = badger.Open(ctx, logsDir)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer lc.store.Close()

	for logName, logState := range lc.state.Logs {
		lc.startWorker(ctx, logName, logState)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				go func() {
					if err := lc.agent.Send(func() error {
						return lc.removeOldEvents(ctx)
					}); err != nil {
						lc.agent.Fail(fmt.Errorf("removing old events: %w", err))
					}
				}()
			}
		}
	}()

	return lc.agent.Run(ctx)
}

func (lc *LogCollector) removeOldEvents(ctx context.Context) error {
	lc.debugf("removing old events")
	for logName := range lc.state.Logs {
		log := lc.store.GetLog(logName)
		if err := log.RemoveOldEvents(ctx); err != nil {
			return fmt.Errorf("removing %q events: %w", logName, err)
		}
	}
	lc.debugf("removed old events")
	return nil
}

func validLogName(s string) bool {
	return s != "" // TODO: More validation. Cannot have internal null bytes.
}

func (lc *LogCollector) AddLog(ctx context.Context, input *api.AddLogInput) (output *api.AddLogOutput, err error) {
	err = lc.agent.Send(func() error {
		output, err = lc.addLog(ctx, input)
		return err
	})
	return
}

func (lc *LogCollector) addLog(ctx context.Context, input *api.AddLogInput) (*api.AddLogOutput, error) {
	if !validLogName(input.Name) {
		return nil, fmt.Errorf("invalid log name: %q", input.Name)
	}
	if input.Source == "" {
		return nil, errors.New("log source path is required")
	}
	if lc.state.Logs == nil {
		lc.state.Logs = make(map[string]*logState)
	}
	if lc.state.Logs[input.Name] != nil {
		return nil, fmt.Errorf("already have log %q", input.Name)
	}
	logState := &logState{
		Source: input.Source,
	}
	lc.state.Logs[input.Name] = logState
	if err := lc.saveState(); err != nil {
		return nil, fmt.Errorf("saving new log source: %w", err)
	}

	lc.startWorker(ctx, input.Name, logState)
	return &api.AddLogOutput{}, nil
}

func (lc *LogCollector) startWorker(ctx context.Context, logName string, state *logState) {
	ctx, stop := context.WithCancel(ctx)
	wkr := &collectorWorker{
		Worker: Worker{
			Source: state.Source,
			Sink:   lc.store.GetLog(logName),
			Debug:  lc.debug,
		},
		stop: stop,
	}
	lc.workers[logName] = wkr

	lc.eg.Go(func() error {
		if err := wkr.Run(ctx); err != nil {
			return fmt.Errorf("worker %q run error: %w", logName, err)
		}
		return nil
	})
}

func (lc *LogCollector) RemoveLog(ctx context.Context, input *api.RemoveLogInput) (output *api.RemoveLogOutput, err error) {
	err = lc.agent.Send(func() error {
		output, err = lc.removeLog(ctx, input)
		return err
	})
	return
}

func (lc *LogCollector) removeLog(ctx context.Context, input *api.RemoveLogInput) (*api.RemoveLogOutput, error) {
	delete(lc.state.Logs, input.Name)
	lc.stopWorker(input.Name)
	if err := lc.saveState(); err != nil {
		return nil, fmt.Errorf("saving log source removal: %w", err)
	}
	return &api.RemoveLogOutput{}, nil
}

func (lc *LogCollector) stopWorker(logName string) {
	wkr := lc.workers[logName]
	if wkr == nil {
		return
	}
	wkr.stop()
}

func (lc *LogCollector) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (output *api.DescribeLogsOutput, err error) {
	err = lc.agent.Send(func() error {
		output, err = lc.describeLogs(ctx, input)
		return err
	})
	return
}

func (lc *LogCollector) describeLogs(ctx context.Context, input *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
	var output api.DescribeLogsOutput
	output.Logs = []api.LogDescription{}
	for name, description := range lc.state.Logs {
		var lastEventAt *string
		if lastEvent := lc.store.GetLog(name).GetLastEvent(ctx); lastEvent != nil {
			lastEventAt = &lastEvent.Timestamp
		}

		output.Logs = append(output.Logs, api.LogDescription{
			Name:        name,
			Source:      description.Source,
			LastEventAt: lastEventAt,
		})
	}
	return &output, nil
}

func (lc *LogCollector) GetEvents(ctx context.Context, input *api.GetEventsInput) (output *api.GetEventsOutput, err error) {
	err = lc.agent.Send(func() error {
		output, err = lc.getEvents(ctx, input)
		return err
	})
	return
}

func (lc *LogCollector) getEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	logs := input.Logs
	if logs == nil {
		for name := range lc.state.Logs {
			logs = append(logs, name)
		}
	}

	limit := defaultLimit
	var direction store.Direction
	if input.Next != nil {
		if input.Prev != nil {
			return nil, errors.New("Only one of prev or next may be specified")
		}
		limit = *input.Next
		direction = store.DirectionForward
	} else if input.Prev != nil {
		limit = *input.Prev
		direction = store.DirectionBackward
	} else {
		// Use default limit, and move forward
		direction = store.DirectionForward
	}

	var cursor *store.Cursor
	var err error
	if input.Cursor == nil {
		if cursor, err = lc.getLatestCursor(ctx, input.Logs); err != nil {
			return nil, fmt.Errorf("finding latest cursor: %w", err)
		}
	} else {
		if cursor, err = store.ParseCursor(*input.Cursor); err != nil {
			return nil, fmt.Errorf("parsing cursor: %w", err)
		}
	}

	// TODO: Merge sort.
	events := []api.Event{}
	for _, logName := range logs {
		log := lc.store.GetLog(logName)

		logCursor := cursor
		if logCursor == nil {
			if lastEvent := log.GetLastEvent(ctx); lastEvent != nil {
				idBytes, err := badger.DecodeID(lastEvent.ID)
				if err != nil {
					return nil, fmt.Errorf("decoding id: %w", err)
				}
				idBytes = binaryutil.IncrementBytes(idBytes)
				logCursor = &store.Cursor{ID: idBytes}
			}
		}

		logEvents, err := log.GetEvents(ctx, logCursor, limit, direction)
		if err != nil {
			return nil, fmt.Errorf("getting %q events: %w", logName, err)
		}
		if logEvents == nil {
			continue
		}
		events = append(events, logEvents...)
	}
	sort.Sort(&eventsSorter{events})

	effectiveLimit := mathutil.IntMin(limit, len(events))
	if direction == store.DirectionForward {
		events = events[0:effectiveLimit]
	} else {
		end := len(events)
		start := end - effectiveLimit
		events = events[start:end]
	}

	prevCursor := cursor
	nextCursor := cursor
	if len(events) > 0 {
		// TODO: Separate encoding from badger code.
		if idBytes, err := badger.DecodeID(events[0].ID); err != nil {
			return nil, err
		} else {
			// We don't care about the error if this is already a (zero-valued) NilCursor
			_ = binaryutil.DecrementBytes(idBytes)
			prevCursor = &store.Cursor{ID: idBytes}
		}

		if idBytes, err := badger.DecodeID(events[len(events)-1].ID); err != nil {
			return nil, err
		} else {
			nextCursor = &store.Cursor{ID: binaryutil.IncrementBytes(idBytes)}
		}
	}

	return &api.GetEventsOutput{
		Events:     events,
		PrevCursor: prevCursor.Serialize(),
		NextCursor: nextCursor.Serialize(),
	}, nil
}

func (lc *LogCollector) getLatestCursor(ctx context.Context, logs []string) (*store.Cursor, error) {
	cursor := store.NilCursor
	for _, logName := range logs {
		log := lc.store.GetLog(logName)
		logCursor := log.GetLastCursor(ctx)
		if logCursor != nil && bytes.Compare(logCursor.ID, cursor.ID) > 0 {
			cursor = logCursor
		}
	}

	return cursor, nil
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
