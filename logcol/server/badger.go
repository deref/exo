package server

import (
	"context"
	"encoding/binary"

	badger "github.com/dgraph-io/badger/v3"
)

type badgerSink struct {
	db      *badger.DB
	logName string
}

func newBadgerSink(db *badger.DB, logName string) *badgerSink {
	return &badgerSink{
		db:      db,
		logName: logName,
	}
}

func (sink *badgerSink) AddEvent(ctx context.Context, sid uint64, timestamp uint64, message []byte) error {
	return sink.db.Update(func(txn *badger.Txn) error {
		key := make([]byte, len(sink.logName)+1+8+8)
		pos := copy(key, []byte(sink.logName))
		pos += 1 // Skip null delimiter
		binary.BigEndian.PutUint64(key[pos:], timestamp)
		pos += 8
		binary.BigEndian.PutUint64(key[pos:], sid)

		return txn.Set(key[:], []byte(message))
	})
}
