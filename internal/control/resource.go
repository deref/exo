package control

import (
	"context"
	"errors"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
	. "github.com/deref/exo/sdk"
)

type ResourceComponentConfig struct {
	RawComponentConfig
	Spec map[string]any `json:"spec"`
}

// Implements component methods in terms of an underlying resource.
type ResourceComponentController struct {
	api.Service
}

func (c *ResourceComponentController) InitializeController(ctx context.Context, svc api.Service) error {
	c.Service = svc
	return nil
}

func (c *ResourceComponentController) OnCreate(ctx context.Context, cfg *ResourceComponentConfig) (err error) {
	var m struct {
		Resource struct {
			ID string
		} `graphql:"createResource(type: $type, model: $model, component: $component)"`
	}
	return api.Mutate(ctx, c.Service, &m, map[string]any{
		"type":      cfg.Type,
		"model":     scalars.JSONObject(cfg.Spec),
		"component": cfg.ID,
	})
}

func (c *ResourceComponentController) Render(ctx context.Context, cfg *ResourceComponentConfig) (children []RenderedComponent, err error) {
	// No-children.
	return nil, nil
}

func (c *ResourceComponentController) Shutdown(ctx context.Context, cfg *ResourceComponentConfig) (err error) {
	return errors.New("TODO: shutdown -> delete resource(s)")
}
