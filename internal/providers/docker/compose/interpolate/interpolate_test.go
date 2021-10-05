package interpolate

import (
	"testing"

	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/stretchr/testify/assert"
)

func TestInterpolate(t *testing.T) {
	expected := map[string]interface{}{
		"foo": uint64(123),
		"bar": "i",
		"array": []interface{}{
			"foo",
			"v",
			"bar",
		},
		"map": map[string]interface{}{
			"key":              "x",
			"$notInterpolated": "xyz",
		},
	}
	env := MapEnvironment(map[string]string{
		"one":  "i",
		"five": "v",
		"ten":  "x",
	})
	var v interface{}
	yamlutil.MustUnmarshalString(`
foo: 123
bar: $one
array:
  - foo
  - $five
  - bar
map:
  key: $ten
  $notInterpolated: xyz
`, &v)
	assert.NoError(t, Interpolate(v, env))
	assert.Equal(t, expected, v)
}
