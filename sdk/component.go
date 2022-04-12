package sdk

import (
	"context"

	. "github.com/deref/exo/internal/scalars"
)

type AComponentController = ComponentController[RawJSON]

type ComponentController[Model any] interface {
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

	Run bool
	// TODO: ParitalEnvironment cue.Value `json:"environment"`
	FullEnvironment map[string]string `json:"fullEnvironment"`

	Resources map[string]ComponentConfigResource `json:"resources"`
}

type ComponentConfigResource struct {
	ID   string  `json:"id"`
	Type string  `json:"type"`
	IRI  *string `json:"iri,omitempty"`
}

type RenderedComponent struct {
	Type        string
	Name        string
	Key         string
	Spec        any
	Environment JSONObject
}

// Utilty struct to embed no-op methods for the common case of components that
// only implement RenderComponent.
type PureComponentController[Model any] struct{}

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
