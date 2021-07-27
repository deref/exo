package store

import (
	"context"

	"github.com/deref/exo/logd/api"
)

type Store interface {
	// Returns the next log after the given argument.
	// If given nil, returns the first log.
	// If there are no more logs, returns nil.
	NextLog(after Log) (Log, error)
	GetLog(name string) Log
}

type Log interface {
	Name() string

	// GetEvents returns a page of events along with cursors for moving forward or backward in the result set.
	// If `cursor` is nil, returns the most recent page of events, which is useful for the UI's default tailing behaviour.
	GetEvents(ctx context.Context, cursor *Cursor, limit int, direction Direction) ([]EventWithCursors, error)

	GetLastCursor(context.Context) (*Cursor, error)
	GetLastEvent(context.Context) (*api.Event, error)
	AddEvent(ctx context.Context, timestamp int64, message []byte) error
	// Remove oldest events beyond capacity limit.
	RemoveOldEvents(context.Context) error
	ClearEvents(context.Context) error
}

type Direction int

const (
	DirectionForward  Direction = 1
	DirectionBackward Direction = -1
)

type EventWithCursors struct {
	Event      api.Event
	PrevCursor Cursor
	NextCursor Cursor
}
