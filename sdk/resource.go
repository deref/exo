package sdk

import (
	"context"
	"errors"

	"github.com/deref/exo/internal/api"
	. "github.com/deref/exo/internal/scalars"
)

type AResourceController = ResourceController[RawJSON]

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

func (ctrl *ResourceComponentController) ComponentUpdated(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	switch len(cfg.Resources) {

	case 0:
		var m struct {
			Resource struct {
				ID string
			} `graphql:"createResource(type: $type, model: $model, component: $component)"`
		}
		return api.Mutate(ctx, ctrl.service, &m, map[string]any{
			"type":      cfg.Type,
			"model":     model,
			"component": cfg.ID,
		})

	case 1:
		panic("TODO: update component resource")

	default:
		panic("TODO: handle transitions, report conflicts")

	}
}

func (c *ResourceComponentController) RenderComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) (children []RenderedComponent, err error) {
	// No children.
	return nil, nil
}

func (ctrl *ResourceComponentController) RefreshComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	return errors.New("TODO: Trigger refresh of resources")
}

func (ctrl *ResourceComponentController) ChildrenUpdated(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	// No-op, since there are no children.
	return nil
}

func (c *ResourceComponentController) ShutdownComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) (err error) {
	return errors.New("TODO: delegate to shutdown of resources")
}

func (ctrl *ResourceComponentController) DeleteComponent(ctx context.Context, cfg *ComponentConfig, model *RawJSON) error {
	return errors.New("TODO: delegate to deletion of resources")
}
