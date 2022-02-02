package sdk

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/deref/exo/internal/util/errutil"
)

var ErrResourceGone = errutil.NewHTTPError(http.StatusGone, "resource gone")

// TODO: Improve validation and error reporting for reflective calls.

// TODO: Rename to ResourceController, or move to a resource-specific package.
type ResourceController struct {
	impl reflect.Value
}

// TODO: Most models are JSON, so this interface is a pain. Instead, check
// for json.Marshaler, then encoding.TextMarshaler, and finally just do normal
// JSON marshaling. Similar for unmarshaling.
// TODO: Should we just _insist_ that models are JSON?
type Model interface {
	UnmarshalModel(ctx context.Context, s string) error
	MarshalModel(ctx context.Context) (string, error)
}

func NewResourceController(impl interface{}) *ResourceController {
	return &ResourceController{
		impl: reflect.ValueOf(impl),
	}
}

func unmarshalModel(ctx context.Context, label string, typ reflect.Type, s string) (Model, error) {
	m := reflect.New(typ.Elem()).Interface().(Model)
	err := m.UnmarshalModel(ctx, s)
	if err != nil {
		err = fmt.Errorf("unmarshaling %s: %w", label, err)
	}
	return m, err
}

func marshalModel(ctx context.Context, typ reflect.Type, m Model) (string, error) {
	s, err := m.MarshalModel(ctx)
	if err != nil {
		err = fmt.Errorf("mashaling updated model: %w", err)
	}
	return s, err
}

func (c *ResourceController) Identify(ctx context.Context, model string) (iri string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Identify")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	iri = res[0].Interface().(string)
	err, _ = res[1].Interface().(error)
	return iri, err
}

func (c *ResourceController) Create(ctx context.Context, model string) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Create")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	err, _ = res[0].Interface().(error)
	if err != nil {
		return
	}
	return unmarshaledModel.MarshalModel(ctx)
}

func (c *ResourceController) Read(ctx context.Context, model string) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Read")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	err, _ = res[0].Interface().(error)
	if err != nil {
		return
	}
	return unmarshaledModel.MarshalModel(ctx)
}

func (c *ResourceController) Update(ctx context.Context, prev string, cur string) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Update")
	unmarshaledPrev, err := unmarshalModel(ctx, "previous model", method.Type().In(1), prev)
	if err != nil {
		return "", err
	}
	unmarshaledCur, err := unmarshalModel(ctx, "current model", method.Type().In(2), cur)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledPrev),
		reflect.ValueOf(unmarshaledCur),
	})
	err, _ = res[0].Interface().(error)
	if err != nil {
		return
	}
	return unmarshaledCur.MarshalModel(ctx)
}

func (c *ResourceController) Delete(ctx context.Context, model string) (err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Delete")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	err, _ = res[0].Interface().(error)
	return err
}
