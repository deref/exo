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
	// GetCursorForTimestamp returns a cursor that points to the first entry at or past `timestamp`.
	// If a negative value is given for timestamp, the cursor is for the most recent entry.
	// When no entries are found, this method returns a nil cursor.
	GetCursorForTimestamp(ctx context.Context, timestamp int64) (*Cursor, error)
	GetEvents(ctx context.Context, cursor *Cursor, limit int) (*EventPage, error)
	GetLastEventAt(context.Context) *string
	AddEvent(ctx context.Context, timestamp int64, message []byte) error
	// Remove oldest events beyond capacity limit.
	RemoveOldEvents(context.Context) error
}

type EventPage struct {
	Events []api.Event
	Cursor Cursor
}
