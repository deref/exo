package sdk

import (
	"context"
	"errors"
	"reflect"
)

type Controller struct {
	impl reflect.Value
}

func NewController(impl any) *Controller {
	return &Controller{
		impl: reflect.ValueOf(impl),
	}
}

// TODO: Improve validation and error reporting for reflective calls.

type ResourceComponentConfig struct {
	ComponentConfig
	Model map[string]any
}

// Implements component methods in terms of an underlying resource.
type ResourceComponentController struct{}

func (c *ResourceComponentController) Initialize(ctx context.Context, cfg *ResourceComponentConfig) (err error) {
	return errors.New("TODO: initialize -> create resource")
}

func (c *ResourceComponentController) Render(ctx context.Context, cfg *ResourceComponentConfig) (children []RenderedComponent, err error) {
	// No-children.
	return nil, nil
}

func (c *ResourceComponentController) Shutdown(ctx context.Context, cfg *ResourceComponentConfig) (err error) {
	return errors.New("TODO: shutdown -> delete resource(s)")
}
