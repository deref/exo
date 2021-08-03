package strutil_test

import (
	"testing"

	"github.com/deref/exo/util/strutil"
	"github.com/stretchr/testify/assert"
)

func TestStringEscaping(t *testing.T) {
	testCases := []struct {
		name         string
		in           string
		specialChars []rune
		escapeChar   rune
		expected     string
	}{
		{
			name:         "no escaping",
			in:           "Hello World",
			specialChars: []rune{'$'},
			escapeChar:   '\\',
			expected:     `Hello World`,
		},

		{
			name:         "single character",
			in:           "Hello$World",
			specialChars: []rune{'$'},
			escapeChar:   '\\',
			expected:     `Hello\$World`,
		},

		{
			name:         "multile characters",
			in:           "%foo%",
			specialChars: []rune{'%'},
			escapeChar:   '>',
			expected:     `>%foo>%`,
		},

		{
			name:         "multiple special chars",
			in:           `"Hi" 'there'`,
			specialChars: []rune{'\'', '"'},
			escapeChar:   '\\',
			expected:     `\"Hi\" \'there\'`,
		},

		{
			name:         "escaping the escape char",
			in:           `it\'s`,
			specialChars: []rune{'\''},
			escapeChar:   '\\',
			expected:     `it\\\'s`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out := strutil.EscapeString(testCase.in, testCase.specialChars, testCase.escapeChar)
			assert.Equal(t, testCase.expected, out)

			unescaped := strutil.UnescapeString(out, testCase.escapeChar)
			assert.Equal(t, testCase.in, unescaped)
		})
	}
}
