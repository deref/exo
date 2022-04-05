package sdk

import (
	"context"

	"cuelang.org/go/cue"
)

type AComponentController = ComponentController[cue.Value]

type ComponentController[Model any] interface {
	// Called after a component is first-created.
	ComponentCreated(context.Context, *ComponentConfig, *Model) error
	// Returns a list of desired child components.
	// Called on each iteration of the reconciliation loop.
	RenderComponent(context.Context, *ComponentConfig, *Model) ([]RenderedComponent, error)
	// Called periodically to read state from underlying resources.
	RefreshComponent(context.Context, *ComponentConfig, *Model) error
	// Called after a component has been changed, but before reconciling children.
	ComponentUpdated(context.Context, *ComponentConfig, *Model) error
	// Called when a batch of one or more children have processed the
	// ComponentUpdated hook.
	ChildrenUpdated(context.Context, *ComponentConfig, *Model) error
	// Perform a blocking graceful shutdown of a component.
	ShutdownComponent(context.Context, *ComponentConfig, *Model) error
	// Deletes any associated external resources.
	DeleteComponent(context.Context, *ComponentConfig, *Model) error
}

// Concrete component configuration.
// Corresponds to $Component in the schema.cue file.
type ComponentConfig struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`

	Spec  cue.Value      `json:"spec"`
	State map[string]any `json:"state"`

	Run         bool
	Environment map[string]string `json:"environment"`
}

type RenderedComponent struct {
	Type string
	Name string
	Key  string
	Spec any
}

// Utilty struct to embed no-op methods for the common case of components that
// only implement RenderComponent.
type PureComponentController[Model any] struct{}

func (ctrl *PureComponentController[Model]) ComponentCreated(context.Context, *ComponentConfig, *Model) error {
	return nil
}
func (ctrl *PureComponentController[Model]) RefreshComponent(context.Context, *ComponentConfig, *Model) error {
	return nil
}
func (ctrl *PureComponentController[Model]) ComponentUpdated(context.Context, *ComponentConfig, *Model) error {
	return nil
}
func (ctrl *PureComponentController[Model]) ChildrenUpdated(context.Context, *ComponentConfig, *Model) error {
	return nil
}
func (ctrl *PureComponentController[Model]) ShutdownComponent(context.Context, *ComponentConfig, *Model) error {
	return nil
}
func (ctrl *PureComponentController[Model]) DeleteComponent(context.Context, *ComponentConfig, *Model) error {
	return nil
}

// Adapts a ComponentController[Model] to AComponentController and wraps
// methods with panic recovery.
type ComponentControllerAdapter[Model any] struct {
	impl any
}

func NewComponentControllerAdapater[Model any](impl ComponentController[Model]) AComponentController {
	return &ComponentControllerAdapter[Model]{
		impl: impl,
	}
}

//func (c *ComponentControllerAdapter[Model]) decodeConfig(cfg ComponentConfig[Model]) (ComponentConfig[Model], error) {
//	res := sdk.ComponentConfig[Model]{
//		RawComponentConfig: cfg,
//	}
//	if err := res.Spec.DecodeCue(cfg.SpecValue); err != nil {
//		return res, fmt.Errorf("decoding spec: %w", err)
//	}
//	return res, nil
//}

func (ctrl *ComponentControllerAdapter[Model]) ComponentCreated(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ComponentControllerAdapter[Model]) RenderComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) ([]RenderedComponent, error) {
	panic("TODO")
}

func (ctrl *ComponentControllerAdapter[Model]) RefreshComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ComponentControllerAdapter[Model]) ComponentUpdated(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ComponentControllerAdapter[Model]) ChildrenUpdated(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ComponentControllerAdapter[Model]) ShutdownComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	panic("TODO")
}

func (ctrl *ComponentControllerAdapter[Model]) DeleteComponent(ctx context.Context, cfg *ComponentConfig, model *cue.Value) error {
	panic("TODO")
}
