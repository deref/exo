package control

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/sdk"
)

// Adapts a controller implementation to the generic controller interface.
// Wrapped methods have automatic spec encoding/decoding and panic handling.
// This is intended for use with trusted, in-process controllers.
type ControllerAdapter[M Model] struct {
	impl sdk.Controller[M]
}

var _ Controller = (*ControllerAdapter[Model])(nil)

func AdaptController[M Model](ctx context.Context, svc api.Service, impl sdk.Controller[M]) (*ControllerAdapter[M], error) {
	c := &ControllerAdapter[M]{
		impl: impl,
	}
	if err := c.InitializeController(ctx, svc); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ControllerAdapter[M]) decodeConfig(cfg ComponentConfig) (sdk.ComponentConfig[M], error) {
	res := sdk.ComponentConfig[M]{
		RawComponentConfig: cfg,
	}
	if err := res.Spec.DecodeCue(cfg.SpecValue); err != nil {
		return res, fmt.Errorf("decoding spec: %w", err)
	}
	return res, nil
}

func (c *ControllerAdapter[M]) InitializeController(ctx context.Context, svc api.Service) (err error) {
	defer errutil.RecoverTo(&err)
	return c.impl.InitializeController(ctx, svc)
}

func (c *ControllerAdapter[M]) InitializeComponent(ctx context.Context, cfg ComponentConfig) (err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return err
	}
	return c.impl.InitializeComponent(ctx, scfg)
}

func (c *ControllerAdapter[M]) ComponentCreated(ctx context.Context, cfg ComponentConfig) (err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return err
	}
	return c.impl.ComponentCreated(ctx, scfg)
}

func (c *ControllerAdapter[M]) RenderComponent(ctx context.Context, cfg ComponentConfig) (children []RenderedComponent, err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return nil, err
	}
	return c.impl.RenderComponent(ctx, scfg)
}

func (c *ControllerAdapter[M]) RefreshComponent(ctx context.Context, cfg ComponentConfig) (err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return err
	}
	return c.impl.RefreshComponent(ctx, scfg)
}

func (c *ControllerAdapter[M]) ComponentUpdated(ctx context.Context, cfg ComponentConfig) (err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return err
	}
	return c.impl.ComponentUpdated(ctx, scfg)
}

func (c *ControllerAdapter[M]) ChildrenUpdated(ctx context.Context, cfg ComponentConfig) (err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return err
	}
	return c.impl.ChildrenUpdated(ctx, scfg)
}

func (c *ControllerAdapter[M]) ShutdownComponent(ctx context.Context, cfg ComponentConfig) (err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return err
	}
	return c.impl.ShutdownComponent(ctx, scfg)
}

func (c *ControllerAdapter[M]) KillComponent(ctx context.Context, cfg ComponentConfig) (err error) {
	defer errutil.RecoverTo(&err)
	scfg, err := c.decodeConfig(cfg)
	if err != nil {
		return err
	}
	return c.impl.KillComponent(ctx, scfg)
}
