package sdk

import (
	"context"
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
)

type ResourceController[Model any] interface {
	IdentifyResource(context.Context, *ResourceConfig, *Model) (string, error)
	CreateResource(context.Context, *ResourceConfig, *Model) error
	ReadResource(context.Context, *ResourceConfig, *Model) error
	UpdateResource(ctx context.Context, cfg *ResourceConfig, prev *Model, next *Model) error
	ShutdownResource(context.Context, *ResourceConfig, *Model) error
	DeleteResource(context.Context, *ResourceConfig, *Model) error
}

type ResourceConfig struct {
	ID    string    `json:"id"`
	Type  string    `json:"type"`
	IRI   string    `json:"iri,omitempty"`
	Model cue.Value `json:"model"`
}

// Adapts a ResourceController[Model] to ResourceController[cue.Value] and wraps
// methods with panic recovery.
type ResourceControllerAdapater[Model any] struct {
	impl ResourceController[Model]
}

func NewResourceControllerAdapater[Model any](impl ResourceController[Model]) ResourceController[cue.Value] {
	return &ResourceControllerAdapater[Model]{
		impl: impl,
	}
}

func (ctrl *ResourceControllerAdapater[Model]) IdentifyResource(ctx context.Context, cfg *ResourceConfig, model *cue.Value) (string, error) {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) CreateResource(ctx context.Context, cfg *ResourceConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) ReadResource(ctx context.Context, cfg *ResourceConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) UpdateResource(ctx context.Context, cfg *ResourceConfig, prev *cue.Value, next *cue.Value) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) ShutdownResource(ctx context.Context, cfg *ResourceConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) DeleteResource(ctx context.Context, cfg *ResourceConfig, model *cue.Value) error {
	panic("TODO")
}

// Extends a ResourceController with ComponentController methods.
type ResourceComponentController struct {
	service api.Service
	ResourceController[cue.Value]
}

func NewResourceComponentController[Model any](svc api.Service, impl ResourceController[Model]) *ResourceComponentController {
	return &ResourceComponentController{
		service:            svc,
		ResourceController: NewResourceControllerAdapater(impl),
	}
}

func (c *ResourceComponentController) ComponentCreated(ctx context.Context, cfg *ComponentConfig, model *cue.Value) (err error) {
	var m struct {
		Resource struct {
			ID string
		} `graphql:"createResource(type: $type, model: $model, component: $component)"`
	}
	var obj map[string]any
	if err := model.Decode(&obj); err != nil {
		return fmt.Errorf("decoding model: %w", err)
	}
	return api.Mutate(ctx, c.service, &m, map[string]any{
		"type":      cfg.Type,
		"model":     scalars.JSONObject(obj),
		"component": cfg.ID,
	})
}

func (c *ResourceComponentController) RenderComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) (children []RenderedComponent, err error) {
	// No children.
	return nil, nil
}

func (ctrl *ResourceComponentController) RefreshComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	return errors.New("TODO: Trigger refresh of resources")
}

func (ctrl *ResourceComponentController) ComponentUpdated(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	return errors.New("TODO: Trigger update/recreate/transition of component, as needed")
}

func (ctrl *ResourceComponentController) ChildrenUpdated(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	// No-op, since there are no children.
	return nil
}

func (c *ResourceComponentController) ShutdownComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) (err error) {
	return errors.New("TODO: delegate to shutdown of resources")
}

func (ctrl *ResourceComponentController) DeleteComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	return errors.New("TODO: delegate to deletion of resources")
}
