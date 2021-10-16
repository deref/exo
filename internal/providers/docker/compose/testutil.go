package compose

import (
	"reflect"
	"strings"
	"testing"

	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func testYAML(t *testing.T, name string, s string, v interface{}) {
	s = strings.TrimSpace(s)
	t.Run("unmarshal_"+name, func(t *testing.T) {
		out := zeroAddr(reflect.TypeOf(v))
		err := yaml.Unmarshal([]byte(s), out.Interface())
		if assert.NoError(t, err) {
			assert.Equal(t, v, out.Elem().Interface())
		}
	})
	// NOTE: This test is inadequate because lots of MarshalYAML implementations
	// will not exercise the inverse of interpolation. This happens because
	// There are checks like `String.Expression != ""` in order to preserve
	// input.
	// TODO: Make the inverse operation explicit, so it can be tested properly.
	t.Run("marshal_"+name, func(t *testing.T) {
		marshaled, err := yamlutil.MarshalString(v)
		if assert.NoError(t, err) {
			marshaled = strings.TrimSpace(marshaled)
			assert.Equal(t, s, marshaled)
		}
	})
}

func zeroAddr(typ reflect.Type) reflect.Value {
	sliceType := reflect.SliceOf(typ)
	return reflect.MakeSlice(sliceType, 1, 1).Index(0).Addr()
}

func assertInterpolated(t *testing.T, env map[string]string, s string, v interface{}) {
	s = strings.TrimSpace(s)
	out := zeroAddr(reflect.TypeOf(v))
	err := yaml.Unmarshal([]byte(s), out.Interface())
	if !assert.NoError(t, err) {
		return
	}
	err = (out.Interface().(Interpolator)).Interpolate(MapEnvironment(env))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, v, out.Elem().Interface())
}
