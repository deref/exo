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
	assert.False(t, IsValidName("x--y"))
}

func TestMangleName(t *testing.T) {
	check := func(input string, expected string) {
		actual := MangleName(input)
		assert.Equal(t, expected, actual)
	}
	check("foo-bar", "foo-bar")
	check("   asdf    &  --  12321 ---", "asdf-12321")
}
