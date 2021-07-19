package badger

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logd/api"
	"github.com/dgraph-io/badger/v3"
	"github.com/oklog/ulid/v2"
)

func (log *Log) GetLastEventAt(ctx context.Context) *string {
	prefix := append([]byte(log.name), 0)
	var timestamp uint64
	err := log.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 1
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		it.Seek(append([]byte(log.name), 255))
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

func (log *Log) GetEvents(ctx context.Context, after string, limit int) ([]api.Event, error) {
	events := make([]api.Event, 0, limit)
	if err := log.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := append([]byte(log.name), 0)
		start := prefix
		if after != "" {
			id, err := ulid.Parse(strings.ToUpper(after))
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
	return events, nil
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
