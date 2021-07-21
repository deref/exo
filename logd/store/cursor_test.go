package store_test

import (
	"testing"

	"github.com/deref/exo/logd/store"
	"github.com/stretchr/testify/assert"
)

func TestRoundtrip(t *testing.T) {
	cursors := []*store.Cursor{
		{
			ID:        []byte("some-id"),
			Direction: store.DirectionForward,
		},
		{
			ID:        []byte("some-id"),
			Direction: store.DirectionReverse,
		},
	}

	for _, cursor := range cursors {
		serialized := cursor.Serialize()
		parsed, err := store.ParseCursor(serialized)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, cursor, parsed)
	}
}
