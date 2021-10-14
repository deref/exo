package compose

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

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
		bs, err := yaml.Marshal(v)
		if assert.NoError(t, err) {
			bs = bytes.TrimSpace(bs)
			assert.Equal(t, s, string(bs))
		}
	})
}

func zeroAddr(typ reflect.Type) reflect.Value {
	sliceType := reflect.SliceOf(typ)
	return reflect.MakeSlice(sliceType, 1, 1).Index(0).Addr()
}
