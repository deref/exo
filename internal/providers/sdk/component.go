package sdk

import (
	"context"
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/util/errutil"
)

// Concrete component configuration.
// Corresponds to $Component in the schema.cue file, but with
type ComponentConfig struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`

	// Intentionally not called "Spec" and "State". Those names are reserved for
	// more strongly typed field in Model structs that embed a ComponentConfig.
	SpecValue  cue.Value              `json:"spec"`
	StateValue map[string]interface{} `json:"state"`

	Run         bool
	Environment map[string]string `json:"environment"`

	Resources map[string]ResourceConfig
}

type ResourceConfig struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	IRI  string `json:"iri,omitempty"`
}

func decodeComponentConfig(ctx context.Context, typ reflect.Type, v cue.Value) (interface{}, error) {
	out := reflect.New(typ.Elem()).Interface()
	err := v.Decode(out)
	if err != nil {
		return nil, fmt.Errorf("decoding component config: %w", err)
	}
	return out, err
}

func (c *Controller) Initialize(ctx context.Context, cfg cue.Value) (err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Initialize")
	if !method.IsValid() {
		// Default behavior is a no-op.
		return nil
	}
	unmarshaledCfg, err := decodeComponentConfig(ctx, method.Type().In(1), cfg)
	if err != nil {
		return err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledCfg),
	})
	err, _ = res[0].Interface().(error)
	return err
}

func (c *Controller) Render(ctx context.Context, cfg cue.Value) (children []RenderedComponent, err error) {
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
	return res[0].Interface().([]RenderedComponent), nil
}

type RenderedComponent struct {
	Type string
	Name string
	Key  string
	Spec interface{}
}
