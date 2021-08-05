package cmdutil_test

import (
	"testing"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected *cmdutil.ParsedArgs
	}{
		{
			name: "empty",
			args: []string{},
			expected: &cmdutil.ParsedArgs{
				Command: "",
				Args:    []string{},
				Flags:   map[string]string{},
			},
		},

		{
			name: "bare command",
			args: []string{"my-cmd"},
			expected: &cmdutil.ParsedArgs{
				Command: "my-cmd",
				Args:    []string{},
				Flags:   map[string]string{},
			},
		},

		{
			name: "positional args",
			args: []string{"foo", "bar", "baz"},
			expected: &cmdutil.ParsedArgs{
				Command: "foo",
				Args:    []string{"bar", "baz"},
				Flags:   map[string]string{},
			},
		},

		{
			name: "flags",
			args: []string{"foo", "--bar", "baz", "-v", "on"},
			expected: &cmdutil.ParsedArgs{
				Command: "foo",
				Args:    []string{},
				Flags: map[string]string{
					"bar": "baz",
					"v":   "on",
				},
			},
		},

		{
			name: "flags and positional args",
			args: []string{"tar", "-xzvf", "foo.tar.gz", "my-dir"},
			expected: &cmdutil.ParsedArgs{
				Command: "tar",
				Args:    []string{"my-dir"},
				Flags: map[string]string{
					"xzvf": "foo.tar.gz",
				},
			},
		},

		{
			name: "flags with embedded value",
			args: []string{"server", "--port=8080"},
			expected: &cmdutil.ParsedArgs{
				Command: "server",
				Args:    []string{},
				Flags: map[string]string{
					"port": "8080",
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := cmdutil.ParseArgs(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, parsed)
		})
	}
}
