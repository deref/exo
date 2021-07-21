package badger

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/logd/store"
	"github.com/dgraph-io/badger/v3"
	"github.com/oklog/ulid/v2"
)

func (log *Log) GetCursorForTimestamp(ctx context.Context, timestamp int64) (*store.Cursor, error) {
	if timestamp < 0 {
		lastEvent, err := log.getLastEvent(ctx)
		if err != nil {
			return nil, err
		}
		if lastEvent == nil {
			return nil, nil
		}

		id := ulid.MustParse(strings.ToUpper(lastEvent.ID))
		return &store.Cursor{
			ID:        id[:],
			Direction: store.DirectionForward,
		}, nil
	}

	logNameLen := len(log.name)
	prefix := append([]byte(log.name), 0)
	var nextID []byte
	var foundNext bool

	if err := log.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		// TODO: Add an index on (log, timestamp) -> id to avoid the table scan.
		for it.Seek(prefix); it.ValidForPrefix(prefix) && !foundNext; it.Next() {
			item := it.Item()
			key := item.Key()

			idOffset := logNameOffset + logNameLen + 1
			nextID = key[idOffset:]

			return item.Value(func(val []byte) error {
				if err := validateVersion(val[versionOffset]); err != nil {
					return err
				}
				eventTimestamp := int64(binary.BigEndian.Uint64(val[timestampOffset : timestampOffset+timestampLen]))
				if timestamp <= eventTimestamp {
					foundNext = true
				}
				return nil
			})
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("scanning logs: %w", err)
	}

	// If no entry found with timestamp at least as new as the one supplied, return nil Cursor.
	if !foundNext {
		return nil, nil
	}

	return &store.Cursor{
		ID:        nextID,
		Direction: store.DirectionForward,
	}, nil
}

func (log *Log) GetLastEventAt(ctx context.Context) *string {
	lastEvent, err := log.getLastEvent(ctx)
	if err != nil || lastEvent == nil {
		return nil
	}
	return &lastEvent.Timestamp
}

func (log *Log) getLastEvent(ctx context.Context) (*api.Event, error) {
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
			key := it.Item().Key()
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

func (log *Log) GetEvents(ctx context.Context, cursor *store.Cursor, limit int) (*store.EventPage, error) {
	prefix := append([]byte(log.name), 0)
	start := prefix
	if cursor != nil {
		if !cursor.IsValid() {
			return nil, errors.New("cursor is not valid")
		}

		if cursor.Direction != store.DirectionForward {
			return nil, errors.New("reverse pagination not yet supported")
		}

		start = append(start, incrementBytes(cursor.ID)...)
	}

	events := make([]api.Event, 0, limit)
	if err := log.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		eventsProcessed := 0
		for it.Seek(start); it.ValidForPrefix(prefix) && eventsProcessed < limit; it.Next() {
			item := it.Item()
			key := item.Key()
			if err := item.Value(func(val []byte) error {
				evt, err := eventFromEntry(log.name, key, val)
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
	page := &store.EventPage{
		Events: events,
	}

	// Set page cursor from last event.
	if len(events) > 0 {
		lastEvent := events[len(events)-1]
		ulid, err := ulid.Parse(strings.ToUpper(lastEvent.ID))
		if err != nil {
			return nil, fmt.Errorf("parsing cursor event ID for cursor: %w", err)
		}
		page.Cursor = store.Cursor{
			ID:        ulid[:],
			Direction: store.DirectionForward,
		}
	}

	return page, nil
}

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
