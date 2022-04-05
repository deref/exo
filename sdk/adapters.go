package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	. "github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/errutil"
)

// Adapts a ComponentController[Model] to AComponentController and wraps
// methods with panic recovery.
type ComponentControllerAdapter[Model any] struct {
	impl ComponentController[Model]
}

func NewComponentControllerAdapater[Model any](impl ComponentController[Model]) AComponentController {
	return &ComponentControllerAdapter[Model]{
		impl: impl,
	}
}

func (ctrl *ComponentControllerAdapter[Model]) RenderComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) ([]RenderedComponent, error) {
	res, err := callAdapted[Model](ctx, ctrl.impl.RenderComponent, cfg, modelAdapter{"model", model})
	children, _ := res.([]RenderedComponent)
	return children, err
}

func (ctrl *ComponentControllerAdapter[Model]) RefreshComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.RefreshComponent, cfg, modelAdapter{"model", model})
	return err
}

func (ctrl *ComponentControllerAdapter[Model]) ComponentUpdated(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.ComponentUpdated, cfg, modelAdapter{"model", model})
	return err
}

func (ctrl *ComponentControllerAdapter[Model]) ChildrenUpdated(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.ChildrenUpdated, cfg, modelAdapter{"model", model})
	return err
}

func (ctrl *ComponentControllerAdapter[Model]) ShutdownComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.ShutdownComponent, cfg, modelAdapter{"model", model})
	return err
}

func (ctrl *ComponentControllerAdapter[Model]) DeleteComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.DeleteComponent, cfg, modelAdapter{"model", model})
	return err
}

// Adapts a ResourceController[Model] to AResourceController and wraps
// methods with panic recovery.
type ResourceControllerAdapater[Model any] struct {
	impl ResourceController[Model]
}

func NewResourceControllerAdapater[Model any](impl ResourceController[Model]) AResourceController {
	return &ResourceControllerAdapater[Model]{
		impl: impl,
	}
}

func (ctrl *ResourceControllerAdapater[Model]) IdentifyResource(ctx context.Context, cfg *ResourceConfig, model *RawJSON) (string, error) {
	res, err := callAdapted[Model](ctx, ctrl.impl.IdentifyResource, cfg, modelAdapter{"model", model})
	iri, _ := res.(string)
	return iri, err
}

func (ctrl *ResourceControllerAdapater[Model]) CreateResource(ctx context.Context, cfg *ResourceConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.CreateResource, cfg, modelAdapter{"model", model})
	return err
}

func (ctrl *ResourceControllerAdapater[Model]) ReadResource(ctx context.Context, cfg *ResourceConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.ReadResource, cfg, modelAdapter{"model", model})
	return err
}

func (ctrl *ResourceControllerAdapater[Model]) UpdateResource(ctx context.Context, cfg *ResourceConfig, prev *RawJSON, next *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.UpdateResource, cfg, modelAdapter{"previous model", prev}, modelAdapter{"next model", next})
	return err
}

func (ctrl *ResourceControllerAdapater[Model]) ShutdownResource(ctx context.Context, cfg *ResourceConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.ShutdownResource, cfg, modelAdapter{"model", model})
	return err
}

func (ctrl *ResourceControllerAdapater[Model]) DeleteResource(ctx context.Context, cfg *ResourceConfig, model *RawJSON) error {
	_, err := callAdapted[Model](ctx, ctrl.impl.DeleteResource, cfg, modelAdapter{"model", model})
	return err
}

type modelAdapter struct {
	Label string
	Model *RawJSON
}

func callAdapted[Model any](ctx context.Context, f any, args ...any) (res any, err error) {
	defer errutil.RecoverTo(&err)

	// Marshal models.
	adaptedArgs := make([]reflect.Value, 1+len(args))
	adaptedArgs[0] = reflect.ValueOf(ctx)
	for i, arg := range args {
		var adaptedArg reflect.Value
		if ma, ok := arg.(modelAdapter); ok {
			var model Model
			if unmarshalErr := json.Unmarshal(*ma.Model, &model); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshaling %s: %w", ma.Label, unmarshalErr)
			}
			adaptedArg = reflect.ValueOf(&model)
		} else {
			adaptedArg = reflect.ValueOf(arg)
		}
		adaptedArgs[1+i] = adaptedArg
	}

	// Invoke adapted method.
	fValue := reflect.ValueOf(f)
	results := fValue.Call(adaptedArgs)

	// Extract (result any, err error) return values.
	// The result return value is optional.
	var rerr reflect.Value
	switch len(results) {
	case 1:
		rerr = results[0]
	case 2:
		res = results[0].Interface()
		rerr = results[1]
	default:
		panic("unreachable")
	}
	if !rerr.IsNil() {
		err = rerr.Interface().(error)
	}

	// Unmarshal updated models.
	for i, arg := range args {
		if ma, ok := arg.(modelAdapter); ok {
			adaptedArg := adaptedArgs[1+i]
			marshaledModel, marshalErr := json.Marshal(adaptedArg.Interface())
			if marshalErr != nil {
				err = fmt.Errorf("unmarshaling %s: %w", ma.Label, marshalErr)
				return
			}
			reflect.ValueOf(ma.Model).Elem().Set(reflect.ValueOf(RawJSON(marshaledModel)))
		}
	}

	return
}
