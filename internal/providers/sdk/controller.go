package sdk

import (
	"context"
	"errors"
	"reflect"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/errutil"
)

type Controller struct {
	impl reflect.Value
}

func NewController(ctx context.Context, svc api.Service, impl any) (*Controller, error) {
	c := &Controller{
		impl: reflect.ValueOf(impl),
	}
	if err := c.Init(ctx, svc); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Controller) Init(ctx context.Context, svc api.Service) (err error) {
	defer errutil.RecoverTo(&err)
	method := c.impl.MethodByName("Init")
	if !method.IsValid() {
		// Default is a no-op.
		return nil
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(svc),
	})
	err, _ = res[0].Interface().(error)
	return err
}

// TODO: Improve validation and error reporting for reflective calls.

type ResourceComponentConfig struct {
	ComponentConfig
	Model map[string]any
}

// Implements component methods in terms of an underlying resource.
type ResourceComponentController struct {
	api.Service
}

func (c *ResourceComponentController) Init(ctx context.Context, svc api.Service) error {
	c.Service = svc
	return nil
}

func (c *ResourceComponentController) OnCreate(ctx context.Context, cfg *ResourceComponentConfig) (err error) {
	return errors.New("TODO: resource component controller onCreate -> create resource")
}

func (c *ResourceComponentController) Render(ctx context.Context, cfg *ResourceComponentConfig) (children []RenderedComponent, err error) {
	// No-children.
	return nil, nil
}

func (c *ResourceComponentController) Shutdown(ctx context.Context, cfg *ResourceComponentConfig) (err error) {
	return errors.New("TODO: shutdown -> delete resource(s)")
}
