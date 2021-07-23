package binaryutil_test

import (
	"testing"

	"github.com/deref/exo/util/binaryutil"
	"github.com/stretchr/testify/assert"
)

func TestIncrementBytes(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{
			input:    nil,
			expected: nil,
		},
		{
			input:    []byte{},
			expected: []byte{1},
		},
		{
			input:    []byte{0},
			expected: []byte{1},
		},
		{
			input:    []byte{1},
			expected: []byte{2},
		},
		{
			input:    []byte{255},
			expected: []byte{1, 0},
		},
	}

	for _, testCase := range testCases {
		out := binaryutil.IncrementBytes(testCase.input)
		assert.Equal(t, testCase.expected, out)
	}
}

func TestDecrementBytes(t *testing.T) {
	testCases := []struct {
		input       []byte
		expected    []byte
		expectError bool
	}{
		{
			input:    nil,
			expected: nil,
		},
		{
			input:       []byte{},
			expectError: true,
		},
		{
			input:       []byte{0},
			expectError: true,
		},
		{
			input:    []byte{1},
			expected: []byte{0},
		},
		{
			input:    []byte{2},
			expected: []byte{1},
		},
		{
			input: []byte{1, 0},
			// Does not reallocate, so there may be leading `0` bytes
			expected: []byte{0, 255},
		},
	}

	for _, testCase := range testCases {
		bs := testCase.input
		err := binaryutil.DecrementBytes(bs)
		if testCase.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, bs)
		}
	}
}
