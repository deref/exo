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
		key := make([]byte, len(sink.logName)+1+8)
		copy(key, []byte(sink.logName))
		binary.BigEndian.PutUint64(key[len(sink.logName)+1:], sid)

		val := make([]byte, 8+len(message))
		binary.BigEndian.PutUint64(val, timestamp)
		copy(val[9:], message)

		return txn.Set(key[:], val)
	})
}
