package procfile

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	assertParse := func(expected Process, command string) {
		actual, err := ParseCommand(strings.NewReader(command))
		assert.NoError(t, err)
		assert.Equal(t, expected, *actual)
	}

	assertParse(Process{
		Program:     "program",
		Arguments:   []string{},
		Environment: map[string]string{},
	}, "program")

	assertParse(Process{
		Program:     "program",
		Arguments:   []string{"arg"},
		Environment: map[string]string{},
	}, "program arg")

	assertParse(Process{
		Program:     "program",
		Arguments:   []string{"quoted arg"},
		Environment: map[string]string{},
	}, `program "quoted arg"`)

	assertParse(Process{
		Program:   "foo",
		Arguments: []string{"bar", "baz"},
		Environment: map[string]string{
			"x": "1",
			"y": "2",
		},
	}, "x=1 y=2 foo bar baz")

	assertParse(Process{
		Program:     "3",
		Arguments:   []string{},
		Environment: map[string]string{},
	}, "$((1+2))")

	assertParse(Process{
		Program:     "axb ayb",
		Arguments:   []string{},
		Environment: map[string]string{},
	}, "a{x,y}b")

	assertNoParse := func(errSubstr string, command string) {
		_, err := ParseCommand(strings.NewReader(command))
		assert.Error(t, err)
		if !strings.Contains(err.Error(), errSubstr) {
			t.Errorf("expected parsing %q to produce error containing %q, got %q", command, errSubstr, err.Error())
		}
	}
	assertNoParse("reached EOF without closing quote", `"`)
	assertNoParse("unsupported: command substitution", `foo $(echo 123)`)
	assertNoParse("unsupported: glob patterns", `foo *`)
	assertNoParse("unbound variable", `$X`)
}
