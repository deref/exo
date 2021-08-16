package compose

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/deref/exo/internal/providers/docker/compose/template"
)

type Interpolator interface {
	Interpolate(Environment) error
}

type Environment = template.Environment

type MapEnvironment = template.MapEnvironment

func Interpolate(v interface{}, env Environment) error {
	var visitor interpolationVisitor
	visitor.visit(v, env)
	return visitor.err
}

type interpolationVisitor struct {
	path []string
	err  error
}

func (iv *interpolationVisitor) setErr(err error) {
	if iv.err == nil && err != nil {
		path := strings.Join(iv.path, "/")
		iv.err = fmt.Errorf("at /%s: %w", path, err)
	}
}

func (iv *interpolationVisitor) visit(v interface{}, env Environment) {
	interpolator, ok := v.(Interpolator)
	if ok {
		iv.setErr(interpolator.Interpolate(env))
		return
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		panic("interpolation of non-pointer")
	}
	rv = rv.Elem()

	switch rv.Kind() {

	case reflect.Struct:
		typ := rv.Type()
		for fieldIndex := 0; fieldIndex < typ.NumField(); fieldIndex++ {
			field := typ.Field(fieldIndex)
			tag := field.Tag.Get("yaml")
			if tag == "" {
				iv.setErr(fmt.Errorf("no yaml tag for field: %s.%s", typ.Name(), field.Name))
				return
			}
			iv.path = append(iv.path, tag)
			iv.visit(rv.Field(fieldIndex).Addr().Interface(), env)
			iv.path = iv.path[:len(iv.path)-1]
		}

	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			iv.visit(iter.Value(), env)
		}

	default:
		iv.setErr(fmt.Errorf("cannot interpolate %s", rv.Type()))
		return
	}
}
