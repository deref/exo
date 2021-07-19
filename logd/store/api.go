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
	GetEvents(ctx context.Context, after string, limit int) ([]api.Event, error)
	GetLastEventAt(context.Context) *string
	AddEvent(ctx context.Context, timestamp uint64, message []byte) error
	// Remove oldest events beyond capacity limit.
	RemoveOldEvents(context.Context) error
}
