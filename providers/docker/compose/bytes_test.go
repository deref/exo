package compose

import (
	"strings"
	"testing"

	"github.com/deref/exo/util/yamlutil"
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
		Int Bytes
		Str Bytes
	}
	yamlutil.MustUnmarshalString(`
int: 1234
str: 5k
`, &actual)
	assert.Equal(t, int64(1234), int64(actual.Int))
	assert.Equal(t, int64(5120), int64(actual.Str))

	assert.Equal(t, "1024", strings.TrimSpace(yamlutil.MustMarshalString(Bytes(1024))))
}
