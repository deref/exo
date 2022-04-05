package sdk

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/deref/exo/internal/api"
)

type AResourceController = ResourceController[json.RawMessage]

type ResourceController[Model any] interface {
	IdentifyResource(context.Context, *ResourceConfig, *Model) (string, error)
	CreateResource(context.Context, *ResourceConfig, *Model) error
	ReadResource(context.Context, *ResourceConfig, *Model) error
	UpdateResource(ctx context.Context, cfg *ResourceConfig, prev *Model, next *Model) error
	ShutdownResource(context.Context, *ResourceConfig, *Model) error
	DeleteResource(context.Context, *ResourceConfig, *Model) error
}

type ResourceConfig struct {
	ID   string  `json:"id"`
	Type string  `json:"type"`
	IRI  *string `json:"iri,omitempty"`
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

func (ctrl *ResourceControllerAdapater[Model]) IdentifyResource(ctx context.Context, cfg *ResourceConfig, model *json.RawMessage) (string, error) {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) CreateResource(ctx context.Context, cfg *ResourceConfig, model *json.RawMessage) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) ReadResource(ctx context.Context, cfg *ResourceConfig, model *json.RawMessage) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) UpdateResource(ctx context.Context, cfg *ResourceConfig, prev *json.RawMessage, next *json.RawMessage) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) ShutdownResource(ctx context.Context, cfg *ResourceConfig, model *json.RawMessage) error {
	panic("TODO")
}

func (ctrl *ResourceControllerAdapater[Model]) DeleteResource(ctx context.Context, cfg *ResourceConfig, model *json.RawMessage) error {
	panic("TODO")
}

// Extends a ResourceController with ComponentController methods.
type ResourceComponentController struct {
	service api.Service
	AResourceController
}

func NewResourceComponentController[Model any](svc api.Service, impl ResourceController[Model]) *ResourceComponentController {
	return &ResourceComponentController{
		service:             svc,
		AResourceController: NewResourceControllerAdapater(impl),
	}
}

func (c *ResourceComponentController) ComponentCreated(ctx context.Context, cfg *ComponentConfig, model *json.RawMessage) (err error) {
	var m struct {
		Resource struct {
			ID string
		} `graphql:"createResource(type: $type, model: $model, component: $component)"`
	}
	return api.Mutate(ctx, c.service, &m, map[string]any{
		"type":      cfg.Type,
		"model":     model,
		"component": cfg.ID,
	})
}

func (c *ResourceComponentController) RenderComponent(ctx context.Context, cfg *ComponentConfig, model *json.RawMessage) (children []RenderedComponent, err error) {
	// No children.
	return nil, nil
}

func (ctrl *ResourceComponentController) RefreshComponent(ctx context.Context, cfg *ComponentConfig, model *json.RawMessage) error {
	return errors.New("TODO: Trigger refresh of resources")
}

func (ctrl *ResourceComponentController) ComponentUpdated(ctx context.Context, cfg *ComponentConfig, model *json.RawMessage) error {
	return errors.New("TODO: Trigger update/recreate/transition of component, as needed")
}

func (ctrl *ResourceComponentController) ChildrenUpdated(ctx context.Context, cfg *ComponentConfig, model *json.RawMessage) error {
	// No-op, since there are no children.
	return nil
}

func (c *ResourceComponentController) ShutdownComponent(ctx context.Context, cfg *ComponentConfig, model *json.RawMessage) (err error) {
	return errors.New("TODO: delegate to shutdown of resources")
}

func (ctrl *ResourceComponentController) DeleteComponent(ctx context.Context, cfg *ComponentConfig, model *json.RawMessage) error {
	return errors.New("TODO: delegate to deletion of resources")
}
