package compose

import (
	"strings"
	"testing"

	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/stretchr/testify/assert"
)

func TestParseBytes(t *testing.T) {
	check := func(s string, expected int64) {
		actual, err := ParseBytes(s)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	check("5", 5)
	check("10b", 10)
	check("2k", 2048)
	check("2kb", 2048)
	check("1m", 1024*1024)
	check("1mb", 1024*1024)
	check("1g", 1024*1024*1024)
	check("1gb", 1024*1024*1024)
}

func TestBytesYaml(t *testing.T) {
	var actual struct {
		Int Bytes `yaml:"int"`
		Str Bytes `yaml:"str"`
		Var Bytes `yaml:"var"`
	}
	mustLoadYaml(`
int: 1234
str: 5k
var: ${ten}k
`, &actual, MapEnvironment{
		"ten": "10",
	})
	assert.Equal(t, int64(1234), int64(actual.Int.Value))
	assert.Equal(t, int64(5120), int64(actual.Str.Value))
	assert.Equal(t, int64(10240), int64(actual.Var.Value))

	assert.Equal(t, `1024`, strings.TrimSpace(yamlutil.MustMarshalString(Bytes{
		Expression: "1024",
	})))
	assert.Equal(t, `1024mb`, strings.TrimSpace(yamlutil.MustMarshalString(Bytes{
		Expression: "1024mb",
	})))
}
