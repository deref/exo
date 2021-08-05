package store_test

import (
	"testing"

	"github.com/deref/exo/internal/logd/store"
	"github.com/stretchr/testify/assert"
)

func TestCursorRoundTrip(t *testing.T) {
	initial := "01fb7m0krevg0kkkkqtd66a0xA"
	cursor, err := store.ParseCursor(initial)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, initial, cursor.Serialize())
}
