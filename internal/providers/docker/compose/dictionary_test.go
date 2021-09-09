package compose

import (
	"strings"
	"testing"

	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/stretchr/testify/assert"
)

func strAddr(s string) *string {
	return &s
}

func TestDictionaryYaml(t *testing.T) {
	type Data struct {
		Dict Dictionary `yaml:"dict"`
	}

	data := Data{
		Dict: Dictionary(map[string]*string{
			"a": strAddr("1"),
			"b": strAddr("2"),
		}),
	}
	mapStr := `
dict:
    a: "1"
    b: "2"
`
	arrayStr := `
dict:
  - a=1
  - b=2
`
	assert.Equal(t,
		strings.TrimSpace(mapStr),
		strings.TrimSpace(yamlutil.MustMarshalString(data)),
	)

	{
		var actual Data
		yamlutil.MustUnmarshalString(mapStr, &actual)
		assert.Equal(t, data, actual)
	}

	{
		var actual Data
		yamlutil.MustUnmarshalString(arrayStr, &actual)
		assert.Equal(t, data, actual)
	}
}
