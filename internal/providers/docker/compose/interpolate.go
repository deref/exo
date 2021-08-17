// This file provides machinery for compose file variable interpolation.
//
// References:
// - https://github.com/compose-spec/compose-spec/blob/b369fe5e02d80b619d14974cd1e64e7eea1b2345/spec.md#interpolation
// - https://github.com/docker/compose/blob/4a51af09d6cdb9407a6717334333900327bc9302/compose/config/interpolation.py
package compose

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/deref/exo/internal/providers/docker/compose/template"
)

// Interpolate mutates v, interpolating contained "Expressions" and storing
// the coercing the results in to corresponding "Value" fields. Some standard
// types are interpolated automatically; types may override their interpolation
// behavior by implementing the Interpolator interface.
func Interpolate(v interface{}, env Environment) error {
	var visitor interpolationVisitor
	visitor.visit(v, env)
	return visitor.err
}

type Interpolator interface {
	Interpolate(Environment) error
}

type Environment = template.Environment

type MapEnvironment = template.MapEnvironment

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
			key, ok := iter.Key().Interface().(string)
			if !ok {
				iv.setErr(fmt.Errorf("expected key string, got: %s", iter.Key().Type()))
				return
			}
			// TODO: Is there a better way to get an addressable item than to
			// create a slice of oen element?
			vslice := reflect.MakeSlice(reflect.SliceOf(iter.Value().Type()), 1, 1)
			velem := vslice.Index(0)
			velem.Set(iter.Value())
			iv.path = append(iv.path, key)
			iv.visit(velem.Addr().Interface(), env)
			rv.SetMapIndex(iter.Key(), velem)
			iv.path = iv.path[:len(iv.path)-1]
		}

	case reflect.Slice:
		n := rv.Len()
		for i := 0; i < n; i++ {
			elem := rv.Index(i)
			iv.path = append(iv.path, strconv.Itoa(i))
			iv.visit(elem.Addr(), env)
			iv.path = iv.path[:len(iv.path)-1]
		}

	case reflect.Ptr:
		if !rv.IsNil() {
			iv.visit(rv.Elem(), env)
		}

	default:
		iv.setErr(fmt.Errorf("cannot interpolate %s", rv.Kind()))
		return
	}
}
