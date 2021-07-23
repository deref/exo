package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	assert.True(t, IsValidName("x"))
	assert.True(t, IsValidName("x123"))
	assert.True(t, IsValidName("x-y"))

	assert.False(t, IsValidName("-"))
	assert.False(t, IsValidName("x-"))
	assert.False(t, IsValidName("5"))
}
