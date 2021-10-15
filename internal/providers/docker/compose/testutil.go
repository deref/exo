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
		sAsV := zeroAddr(reflect.TypeOf(v))
		err := yaml.Unmarshal([]byte(s), sAsV.Interface())
		if assert.NoError(t, err) {
			assert.Equal(t, v, sAsV.Elem().Interface())
		}
	})
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

func boolRef(b bool) *bool {
	return &b
}

func int64Ref(i int64) *int64 {
	return &i
}
