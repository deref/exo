package sdk

import (
	"reflect"
)

type Controller struct {
	impl reflect.Value
}

func NewController(impl interface{}) *Controller {
	return &Controller{
		impl: reflect.ValueOf(impl),
	}
}

// TODO: Improve validation and error reporting for reflective calls.
