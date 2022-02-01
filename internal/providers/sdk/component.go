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

// Concrete component configuration.
// Corresponds to $Component in the schema.cue file.
type ComponentConfig struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Environment map[string]string `json:"environment"`
	// Intentionally not called "Spec". That name is reserved for a more strongly
	// typed field in structs that embed a ComponentConfig.
	SpecValue cue.Value `json:"spec"`
}

func NewComponentController(impl interface{}) *ComponentController {
	return &ComponentController{
		impl: reflect.ValueOf(impl),
	}
}

func (c *ComponentController) Render(ctx context.Context, cfg cue.Value) (_ *RenderResult, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Render")
	unmarshaledCfg, err := decodeComponentConfig(ctx, method.Type().In(1), cfg)
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledCfg),
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

func decodeComponentConfig(ctx context.Context, typ reflect.Type, v cue.Value) (interface{}, error) {
	out := reflect.New(typ.Elem()).Interface()
	err := v.Decode(out)
	if err != nil {
		return nil, fmt.Errorf("decoding component config: %w", err)
	}
	return out, err
}
