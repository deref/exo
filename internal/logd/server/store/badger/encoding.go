package badger

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/logd/api"
	"github.com/oklog/ulid/v2"
)

const (
	eventVersion uint8 = 1

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

func mustIDFromKey(log string, key []byte) string {
	id, err := idFromKey(log, key)
	if err != nil {
		panic(err)
	}
	return id
}

func logFromKey(key []byte) (string, error) {
	// Slice key at separator, and take the first segment as the log name.
	for idx, b := range key {
		if b == 255 {
			return string(key[:idx]), nil
		}
	}

	return "", errors.New("No separator in key")
}

func mustLogFromKey(key []byte) string {
	log, err := logFromKey(key)
	if err != nil {
		panic(err)
	}
	return log
}

func validateVersion(version byte) error {
	if version == 1 {
		return nil
	}
	return fmt.Errorf("unsupported event version: %d; database may have been written with a newer version of exo.", version)
}

func parseID(id []byte) (string, error) {
	var asULID ulid.ULID
	if copy(asULID[:], id) != 16 {
		return "", errors.New("invalid length")
	}

	return strings.ToLower(asULID.String()), nil
}

func decodeID(id string) ([]byte, error) {
	asULID, err := ulid.Parse(strings.ToUpper(id))
	if err != nil {
		return nil, fmt.Errorf("decoding id: %d", err)
	}

	return asULID[:], nil
}
