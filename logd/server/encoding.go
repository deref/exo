package server

import (
	"encoding/binary"
	"fmt"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logd/api"
)

const (
	// Key offsets.
	logNameOffset = 0

	// Value offsets.
	versionOffset   = 0
	versionLen      = 1
	timestampOffset = versionOffset + versionLen
	timestampLen    = 8
	messageOffset   = timestampOffset + timestampLen
)

func eventFromEntry(log string, key, val []byte) (api.Event, error) {
	// Parse key as (logName, null, id).
	id, err := idFromKey(log, key)
	if err != nil {
		return api.Event{}, err
	}

	// Create value as (version, timestamp, message).
	// Version is used so that we can change the value format without rebuilding the database.
	if err := validateVersion(val[versionOffset]); err != nil {
		return api.Event{}, err
	}

	tsNano := binary.BigEndian.Uint64(val[timestampOffset : timestampOffset+timestampLen])
	message := string(val[messageOffset:])

	return api.Event{
		ID:        id,
		Log:       log,
		Timestamp: chrono.NanoToIso(int64(tsNano)),
		Message:   message,
	}, nil
}

func idFromKey(log string, key []byte) (id string, err error) {
	// Parse key as (logName, null, id).
	logNameLen := len(log)
	idOffset := logNameOffset + logNameLen + 1
	id, err = parseID(key[idOffset:])
	if err != nil {
		return "", fmt.Errorf("parsing id: %w", err)
	}
	return id, nil
}

func mustIDFromKey(log string, key []byte) (id string) {
	var err error
	id, err = idFromKey(log, key)
	if err != nil {
		panic(err)
	}
	return id
}

// incrementBytes returns a byte slice that is incremented by 1 bit.
// If `val` is not already only 255-valued bytes, then it is mutated and returned.
// Otherwise, a new slice is allocated and returned.
func incrementBytes(val []byte) []byte {
	for idx := len(val) - 1; idx >= 0; idx-- {
		byt := val[idx]
		if byt == 255 {
			val[idx] = 0
		} else {
			val[idx] = byt + 1
			return val
		}
	}

	// Still carrying from previously most significant byte, so add a new 1-valued byte.
	newVal := make([]byte, len(val)+1)
	newVal[0] = 1
	return newVal
}

func validateVersion(version byte) error {
	if version == 1 {
		return nil
	}
	return fmt.Errorf("unsupported event version: %d; database may have been written with a newer version of exo.", version)
}
