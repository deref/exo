package badger

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/deref/exo/internal/logd/api"
	"github.com/deref/exo/internal/logd/server/store"
	"github.com/deref/exo/internal/util/binaryutil"
	"github.com/dgraph-io/badger/v3"
)

func (log *Log) GetLastEvent(ctx context.Context) (*api.Event, error) {
	var event *api.Event
	prefix := append([]byte(log.name), 0)

	if err := log.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 1
		opts.Reverse = true

		it := txn.NewIterator(opts)
		defer it.Close()

		it.Seek(append([]byte(log.name), 255))
		if it.ValidForPrefix(prefix) {
			item := it.Item()
			key := item.Key()
			return item.Value(func(val []byte) error {
				if err := validateVersion(val[versionOffset]); err != nil {
					return err
				}
				lastEvent, err := eventFromEntry(log.name, key, val)
				if err != nil {
					return err
				}
				event = &lastEvent
				return nil
			})
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("scanning log: %w", err)
	}

	return event, nil
}

func (log *Log) GetLastCursor(ctx context.Context) (*store.Cursor, error) {
	var cursor *store.Cursor
	prefix := append([]byte(log.name), 0)

	err := log.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true

		it := txn.NewIterator(opts)
		defer it.Close()

		it.Seek(append([]byte(log.name), 255))
		if it.ValidForPrefix(prefix) {
			id, err := idFromKey(log.name, it.Item().Key())
			if err != nil {
				return err
			}
			cursor = new(store.Cursor)
			// If `idFromKey` succeeded, this cannot fail.
			cursor.ID, _ = decodeID(id)
			cursor.ID = binaryutil.IncrementBytes(cursor.ID)
		}
		return nil
	})

	return cursor, err
}

func (log *Log) GetEvents(ctx context.Context, cursor *store.Cursor, limit int, direction store.Direction, filterStr *string) ([]store.EventWithCursors, error) {
	prefix := append([]byte(log.name), 0)
	start := prefix
	if cursor != nil {
		start = append(start, cursor.ID...)
	}

	events := make([]api.Event, limit)
	eventsProcessed := 0
	var curIndex int
	var indexDelta int
	if direction == store.DirectionForward {
		curIndex = 0
		indexDelta = 1
	} else {
		curIndex = limit - 1
		indexDelta = -1
	}

	filterStrNorm := ""
	if filterStr != nil {
		filterStrNorm = strings.ToLower(*filterStr)
	}

	if err := log.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		if direction == store.DirectionBackward {
			opts.Reverse = true
		}

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(start); it.ValidForPrefix(prefix) && eventsProcessed < limit; it.Next() {
			item := it.Item()
			key := item.Key()
			if err := item.Value(func(val []byte) error {
				evt, err := eventFromEntry(log.name, key, val)
				if err != nil {
					return err
				}

				// Skip messages that do not match the filter.
				if filterStrNorm != "" && !strings.Contains(strings.ToLower(evt.Message), filterStrNorm) {
					return nil
				}

				events[curIndex] = evt
				curIndex += indexDelta
				eventsProcessed++
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("scanning index: %w", err)
	}

	if direction == store.DirectionForward {
		events = events[0:curIndex]
	} else {
		start := curIndex - indexDelta
		events = events[start:]
	}

	eventsWithCursors := make([]store.EventWithCursors, len(events))
	var err error
	for i, event := range events {
		if eventsWithCursors[i], err = getEventWithCursors(event); err != nil {
			return nil, fmt.Errorf("creating cursors for event %q: %w", event.ID, err)
		}
	}

	return eventsWithCursors, nil
}

func getEventWithCursors(event api.Event) (store.EventWithCursors, error) {
	eventWithCursors := store.EventWithCursors{
		Event: event,
	}
	idBytes, err := decodeID(event.ID)
	if err != nil {
		return store.EventWithCursors{}, err
	}

	// Copy the ID since increment/decrement operations mutate.
	prevCursorID := make([]byte, len(idBytes))
	copy(prevCursorID, idBytes)

	// We don't care about the error if this is already a (zero-valued) NilCursor
	_ = binaryutil.DecrementBytes(prevCursorID)
	eventWithCursors.PrevCursor = store.Cursor{ID: prevCursorID}
	eventWithCursors.NextCursor = store.Cursor{ID: binaryutil.IncrementBytes(idBytes)}

	return eventWithCursors, nil
}

// TODO: RemoveOldEvents also needs to delete events over a certain age,
// not just a count limit. This is necessary for cleanup of no-longer-tracked
// log streams.
const maxEventsPerStream = 5000

func (log *Log) RemoveOldEvents(ctx context.Context) error {
	prefix := append([]byte(log.name), 0)
	var deleteFrom []byte
	if err := log.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		remaining := maxEventsPerStream
		for it.Seek(append([]byte(log.name), 255)); it.ValidForPrefix(prefix); it.Next() {
			if remaining == 0 {
				deleteFrom = it.Item().Key()
				break
			}
			remaining--
		}
		return nil
	}); err != nil {
		return err
	}
	if deleteFrom == nil {
		return nil
	}
	return log.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(deleteFrom); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			if err := txn.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (log *Log) ClearEvents(ctx context.Context) error {
	return errors.New("not implemented: badger.Log.ClearEvents")
}
