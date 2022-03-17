package server

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func Send(ctx context.Context, self any, input any) (output any, err error) {
	selfV := reflect.ValueOf(self)
	inV := reflect.ValueOf(input)
	inT := inV.Type()
	if inT.Kind() != reflect.Ptr {
		panic(fmt.Errorf("expected input to be a pointer, got %s", inT.Kind()))
	}
	inT = inT.Elem()
	methodName := strings.TrimSuffix(inT.Name(), "Input")
	if methodName == inT.Name() {
		panic(fmt.Errorf("expected Input structure, got %T", input))
	}
	method := selfV.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("no method %q on %T", methodName, self)
	}
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		inV,
	})
	if len(results) != 2 {
		panic("expected 2 results")
	}
	output = results[0].Interface()
	errV := results[1]
	if !errV.IsNil() {
		err = errV.Interface().(error)
	}
	return
}
