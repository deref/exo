package compose

import (
	"strings"
	"testing"

	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/stretchr/testify/assert"
)

func TestCommandYaml(t *testing.T) {
	var actual struct {
		Shell  Command
		Parsed Command
	}
	yamlutil.MustUnmarshalString(`
shell: 'foo "bar baz"'
parsed: ['x', 'y z']
`, &actual)
	assert.Equal(t, Command{Parts: []string{`foo "bar baz"`}, IsShellForm: true}, actual.Shell)
	assert.Equal(t, Command{Parts: []string{"x", "y z"}}, actual.Parsed)

	assert.Equal(t, "- x\n- \"y\"", strings.TrimSpace(yamlutil.MustMarshalString(Command{Parts: []string{"x", "y"}})))
}
