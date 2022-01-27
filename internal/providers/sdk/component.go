package sdk

import (
	"context"
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/util/errutil"
)

type ComponentController struct {
	impl reflect.Value
}

func NewComponentController(impl interface{}) *ComponentController {
	return &ComponentController{
		impl: reflect.ValueOf(impl),
	}
}

func unmarshalComponentConfig(ctx context.Context, typ reflect.Type, v cue.Value) (interface{}, error) {
	out := reflect.New(typ.Elem()).Interface()
	err := v.Decode(out)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling spec: %w", err)
	}
	return out, err
}

func (c *ComponentController) Render(ctx context.Context, cfg cue.Value) (_ *RenderResult, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Render")
	unmarshaledSpec, err := unmarshalComponentConfig(ctx, method.Type().In(1), cfg)
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledSpec),
	})
	err, _ = res[1].Interface().(error)
	if err != nil {
		return nil, err
	}
	return res[0].Interface().(*RenderResult), nil
}

type RenderResult struct {
	// resources
	// tasks
}
