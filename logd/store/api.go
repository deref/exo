package store

import (
	"context"
	"io"

	"github.com/deref/exo/logd/api"
)

type Store interface {
	GetLog(name string) Log
	io.Closer // TODO: Remove me.
}

type Log interface {
	// GetEvents returns a page of events along with cursors for moving forward or backward in the result set.
	// If `cursor` is nil, returns the most recent page of events, which is useful for the UI's default tailing behaviour.
	GetEvents(ctx context.Context, cursor *Cursor, limit int, direction Direction) ([]api.Event, error)

	GetLastCursor(context.Context) *Cursor
	GetLastEvent(context.Context) *api.Event
	AddEvent(ctx context.Context, timestamp int64, message []byte) error
	// Remove oldest events beyond capacity limit.
	RemoveOldEvents(context.Context) error
}

type Direction int

const (
	DirectionForward  Direction = 1
	DirectionBackward Direction = -1
)
