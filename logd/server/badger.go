package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/deref/exo/gensym"
	badger "github.com/dgraph-io/badger/v3"
	ulid "github.com/oklog/ulid/v2"
)

const (
	eventVersion uint8 = 1
)

type badgerSink struct {
	db      *badger.DB
	idGen   *gensym.ULIDGenerator
	logName string
}

func newBadgerSink(db *badger.DB, idGen *gensym.ULIDGenerator, logName string) *badgerSink {
	return &badgerSink{
		db:      db,
		idGen:   idGen,
		logName: logName,
	}
}

func (sink *badgerSink) AddEvent(ctx context.Context, timestamp uint64, message []byte) error {
	// Generate an id that is guaranteed to be monotonically increasing within this process.
	id, err := sink.idGen.NextID(ctx)
	if err != nil {
		return fmt.Errorf("generating id: %w", err)
	}

	// Create key as (logName, null, id).
	logNameLen := len(sink.logName)
	idOffset := logNameOffset + logNameLen + 1 // +1 is for null terminator.
	idLen := len(id)                           // Assume that ID the trailing segment and does not need a terminator.

	key := make([]byte, idOffset+idLen)
	copy(key[logNameOffset:logNameOffset+logNameLen], []byte(sink.logName))
	copy(key[idOffset:idOffset+idLen], id)

	// Create value as (version, timestamp, message).
	// Version is used so that we can change the value format without rebuilding the database.
	messageLen := len(message)
	val := make([]byte, messageOffset+messageLen)

	val[versionOffset] = eventVersion
	binary.BigEndian.PutUint64(val[timestampOffset:timestampOffset+timestampLen], timestamp)
	copy(val[messageOffset:messageOffset+messageLen], message)

	return sink.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key[:], val)
	})
}

const maxEventsPerStream = 5000

func (sink *badgerSink) GC(ctx context.Context) error {
	prefix := append([]byte(sink.logName), 0)
	var deleteFrom []byte
	if err := sink.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		remaining := maxEventsPerStream
		for it.Seek(append([]byte(sink.logName), 255)); it.ValidForPrefix(prefix); it.Next() {
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
	return sink.db.Update(func(txn *badger.Txn) error {
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

func parseID(id []byte) (string, error) {
	var asULID ulid.ULID
	if copy(asULID[:], id) != 16 {
		return "", errors.New("invalid length")
	}

	return strings.ToLower(asULID.String()), nil
}
