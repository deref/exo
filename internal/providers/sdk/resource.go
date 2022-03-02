package sdk

import (
	"context"
	"fmt"
	"reflect"

	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/jsonutil"
)

func unmarshalModel(ctx context.Context, label string, typ reflect.Type, obj scalars.JSONObject) (interface{}, error) {
	m := reflect.New(typ.Elem()).Interface()
	err := scalars.DecodeStruct(obj, m)
	if err != nil {
		err = fmt.Errorf("unmarshaling %s: %w", label, err)
	}
	return m, err
}

func marshalModel(ctx context.Context, typ reflect.Type, m interface{}) (string, error) {
	s, err := jsonutil.MarshalString(m)
	if err != nil {
		err = fmt.Errorf("mashaling updated model: %w", err)
	}
	return s, err
}

func (c *Controller) Identify(ctx context.Context, model scalars.JSONObject) (iri string, err error) {
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

func (c *Controller) Create(ctx context.Context, model scalars.JSONObject) (updatedModel string, err error) {
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
		return "", err
	}
	return jsonutil.MarshalString(unmarshaledModel)
}

func (c *Controller) Read(ctx context.Context, model scalars.JSONObject) (updatedModel string, err error) {
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
		return "", err
	}
	return jsonutil.MarshalString(unmarshaledModel)
}

func (c *Controller) Update(ctx context.Context, prev scalars.JSONObject, cur scalars.JSONObject) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Update")
	if !method.IsValid() {
		return "", ErrMethodNotAllowed
	}
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
		return "", err
	}
	return jsonutil.MarshalString(unmarshaledCur)
}

func (c *Controller) Shutdown(ctx context.Context, model scalars.JSONObject) (err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Shutdown")
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

func (c *Controller) Delete(ctx context.Context, model scalars.JSONObject) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Delete")
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
		return "", err
	}
	return jsonutil.MarshalString(unmarshaledModel)
}
