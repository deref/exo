package server

import (
	"context"
	"encoding/binary"
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	eventVersion uint8 = 1
)

type badgerSink struct {
	db      *badger.DB
	idGen   *idGen
	logName string
}

func newBadgerSink(db *badger.DB, idGen *idGen, logName string) *badgerSink {
	return &badgerSink{
		db:      db,
		idGen:   idGen,
		logName: logName,
	}
}

func (sink *badgerSink) AddEvent(ctx context.Context, timestamp uint64, message []byte) error {
	// Generate an id that is guaranteed to be monotonically increasing within this process.
	id, err := sink.idGen.nextId(ctx) // Assume that ID the trailing segment and does not need a terminator.
	if err != nil {
		return fmt.Errorf("generating id: %w", err)
	}

	// Create key as (logName, null, id).
	logNameOffset := 0
	logNameLen := len(sink.logName)
	idOffset := logNameOffset + logNameLen + 1 // +1 is for null terminator.
	idLen := len(id)

	key := make([]byte, idOffset+idLen)
	copy(key[logNameOffset:logNameOffset+logNameLen], []byte(sink.logName))
	copy(key[idOffset:idOffset+idLen], id)

	// Create value as (version, timestamp, message).
	// Version is used so that we can change the value format without rebuilding the database.
	versionOffset := 0
	versionLen := 1
	timestampOffset := versionOffset + versionLen
	timestampLen := 8
	messageOffset := timestampOffset + timestampLen
	messageLen := len(message)
	val := make([]byte, messageOffset+messageLen)

	val[versionOffset] = eventVersion
	binary.BigEndian.PutUint64(val[timestampOffset:timestampOffset+timestampLen], timestamp)
	copy(val[messageOffset:messageOffset+messageLen], message)

	return sink.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key[:], val)
	})
}
