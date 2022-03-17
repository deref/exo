package compose

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/deref/exo/internal/providers/docker/compose/template"
)

type Environment = template.Environment
type MapEnvironment = template.MapEnvironment

var ErrEnvironment = template.ErrEnvironment

type Interpolator interface {
	Interpolate(Environment) error
}

func interpolateStruct(v any, env Environment) error {
	strct := reflect.ValueOf(v).Elem()
	structType := strct.Type()
	numField := structType.NumField()
	for i := 0; i < numField; i++ {
		field := structType.Field(i)
		tag := strings.Split(field.Tag.Get("yaml"), ",")[0]
		if tag == "" {
			panic(fmt.Errorf("expected yaml tag on field: %s", field.Name))
		}
		if tag == "-" {
			continue
		}
		value := strct.Field(i)
		var err error
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			} else {
				value = value.Elem()
			}
		}
		if interpolator, ok := value.Addr().Interface().(Interpolator); ok {
			err = interpolator.Interpolate(env)
		} else if value.Kind() == reflect.Slice {
			err = interpolateSlice(value.Interface(), env)
		} else {
			panic(fmt.Errorf("cannot interpolate field %s of type %s", field.Name, field.Type))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func interpolateSlice(v any, env Environment) error {
	slice := reflect.ValueOf(v)
	n := slice.Len()
	for i := 0; i < n; i++ {
		if err := slice.Index(i).Addr().Interface().(Interpolator).Interpolate(env); err != nil {
			return err
		}
	}
	return nil
}
