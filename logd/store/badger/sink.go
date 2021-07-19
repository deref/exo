package badger

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

func (log *Log) AddEvent(ctx context.Context, timestamp uint64, message []byte) error {
	// Generate an id that is guaranteed to be monotonically increasing within this process.
	id, err := log.idGen.NextID(ctx)
	if err != nil {
		return fmt.Errorf("generating id: %w", err)
	}

	// Create key as (name, null, id).
	logNameLen := len(log.name)
	idOffset := logNameOffset + logNameLen + 1 // +1 is for null terminator.
	idLen := len(id)                           // Assume that ID the trailing segment and does not need a terminator.

	key := make([]byte, idOffset+idLen)
	copy(key[logNameOffset:logNameOffset+logNameLen], []byte(log.name))
	copy(key[idOffset:idOffset+idLen], id)

	// Create value as (version, timestamp, message).
	// Version is used so that we can change the value format without rebuilding the database.
	messageLen := len(message)
	val := make([]byte, messageOffset+messageLen)

	val[versionOffset] = eventVersion
	binary.BigEndian.PutUint64(val[timestampOffset:timestampOffset+timestampLen], timestamp)
	copy(val[messageOffset:messageOffset+messageLen], message)

	return log.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key[:], val)
	})
}
