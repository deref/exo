package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/logd/api"
	"github.com/deref/exo/internal/logd/server/store"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/mathutil"
	"github.com/oklog/ulid/v2"
)

const (
	defaultNextLimit = 500
	maxLimit         = 10000
)

type LogCollector struct {
	Debug bool
	IDGen *gensym.ULIDGenerator
	Store store.Store
}

func (lc *LogCollector) debugf(format string, v ...interface{}) {
	if lc.Debug {
		fmt.Fprintln(os.Stderr, "collector", fmt.Errorf(format, v...))
	}
}

func (lc *LogCollector) AddEvent(ctx context.Context, input *api.AddEventInput) (*api.AddEventOutput, error) {
	if input.Log == "" {
		return nil, errutil.NewHTTPError(http.StatusBadRequest, "log is required")
	}
	log := lc.Store.GetLog(input.Log)
	timestamp, err := chrono.ParseIsoToNano(input.Timestamp)
	if err != nil {
		return nil, errutil.HTTPErrorf(http.StatusBadRequest, "invalid timestamp: %w", err)
	}
	message := []byte(input.Message)
	if err := log.AddEvent(ctx, timestamp, message); err != nil {
		return nil, err
	}
	return &api.AddEventOutput{}, nil
}

func (lc *LogCollector) RemoveOldEvents(ctx context.Context, input *api.RemoveOldEventsInput) (*api.RemoveOldEventsOutput, error) {
	lc.debugf("removing old events")
	var log store.Log
	for {
		var err error
		log, err = lc.Store.NextLog(log)
		if err != nil {
			return nil, fmt.Errorf("enumerating logs: %w", err)
		}
		if log == nil {
			break
		}
		if err := log.RemoveOldEvents(ctx); err != nil {
			return nil, fmt.Errorf("removing %q events: %w", log.Name(), err)
		}
	}
	lc.debugf("removed old events")
	return &api.RemoveOldEventsOutput{}, nil
}

func validLogName(s string) bool {
	return s != "" // TODO: More validation. Cannot have internal null bytes.
}

func (lc *LogCollector) ClearEvents(ctx context.Context, input *api.ClearEventsInput) (output *api.ClearEventsOutput, err error) {
	for _, logName := range input.Logs {
		log := lc.Store.GetLog(logName)
		if err := log.ClearEvents(ctx); err != nil {
			return nil, fmt.Errorf("log %q: %w", logName, err)
		}
	}
	return &api.ClearEventsOutput{}, nil
}

func (lc *LogCollector) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
	var output api.DescribeLogsOutput
	output.Logs = []api.LogDescription{}
	for _, name := range input.Names {
		var lastEventAt *string
		lastEvent, err := lc.Store.GetLog(name).GetLastEvent(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting last event: %w", err)
		}
		if lastEvent != nil {
			lastEventAt = &lastEvent.Timestamp
		}

		output.Logs = append(output.Logs, api.LogDescription{
			Name:        name,
			LastEventAt: lastEventAt,
		})
	}
	return &output, nil
}

func (lc *LogCollector) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	limit := defaultNextLimit
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
		// Use default limit, and move forward.
		direction = store.DirectionForward
	}
	limit = mathutil.IntClamp(limit, 0, maxLimit)

	var cursor *store.Cursor
	var err error
	if input.Cursor == nil {
		if cursor, err = lc.getLatestCursor(ctx, input.Logs); err != nil {
			return nil, fmt.Errorf("finding latest cursor: %w", err)
		}
	} else {
		parsedCursor, err := store.ParseCursor(*input.Cursor)
		if err != nil {
			return nil, fmt.Errorf("parsing cursor: %w", err)
		}
		cursor = &parsedCursor
	}

	// TODO: Merge sort.
	eventsWithCursors := make([]store.EventWithCursors, 0, limit)
	for _, logName := range input.Logs {
		log := lc.Store.GetLog(logName)

		logEventsWithCursors, err := log.GetEvents(ctx, cursor, limit, direction, input.FilterStr)
		if err != nil {
			return nil, fmt.Errorf("getting %q events: %w", logName, err)
		}
		if logEventsWithCursors == nil {
			continue
		}
		eventsWithCursors = append(eventsWithCursors, logEventsWithCursors...)
	}
	sort.Sort(&eventWithCursorsSorter{eventsWithCursors})

	effectiveLimit := mathutil.IntMin(limit, len(eventsWithCursors))
	if direction == store.DirectionForward {
		eventsWithCursors = eventsWithCursors[0:effectiveLimit]
	} else {
		end := len(eventsWithCursors)
		start := end - effectiveLimit
		eventsWithCursors = eventsWithCursors[start:end]
	}

	prevCursor := cursor
	nextCursor := cursor
	if len(eventsWithCursors) > 0 {
		prevCursor = &eventsWithCursors[0].PrevCursor
		nextCursor = &eventsWithCursors[len(eventsWithCursors)-1].NextCursor
	}

	events := make([]api.Event, len(eventsWithCursors))
	for i, eventWithCursors := range eventsWithCursors {
		events[i] = eventWithCursors.Event
	}

	return &api.GetEventsOutput{
		Items:      events,
		PrevCursor: prevCursor.Serialize(),
		NextCursor: nextCursor.Serialize(),
	}, nil
}

func (lc *LogCollector) getLatestCursor(ctx context.Context, logs []string) (*store.Cursor, error) {
	cursor := store.NilCursor

	for _, logName := range logs {
		log := lc.Store.GetLog(logName)
		logCursor, err := log.GetLastCursor(ctx)
		if err != nil {
			return nil, err
		}

		if logCursor != nil && bytes.Compare(logCursor.ID, cursor.ID) > 0 {
			cursor = logCursor
		}
	}

	return cursor, nil
}

type eventWithCursorsSorter struct {
	items []store.EventWithCursors
}

func (iface *eventWithCursorsSorter) Len() int {
	return len(iface.items)
}

func (iface *eventWithCursorsSorter) Less(i, j int) bool {
	iId := ulid.MustParse(iface.items[i].Event.ID)
	jId := ulid.MustParse(iface.items[j].Event.ID)
	return iId.Compare(jId) < 0
}

func (iface *eventWithCursorsSorter) Swap(i, j int) {
	tmp := iface.items[i]
	iface.items[i] = iface.items[j]
	iface.items[j] = tmp
}
