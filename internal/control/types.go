package control

import (
	"context"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/sdk"
)

type ComponentConfig = sdk.RawComponentConfig
type RenderedComponent = sdk.RenderedComponent
type Model = sdk.Model

type Controller interface {
	InitializeController(context.Context, api.Service) error

	InitializeComponent(context.Context, ComponentConfig) error
	ComponentCreated(context.Context, ComponentConfig) error
	RenderComponent(context.Context, ComponentConfig) ([]RenderedComponent, error)
	RefreshComponent(context.Context, ComponentConfig) error
	ComponentUpdated(context.Context, ComponentConfig) error
	ChildrenUpdated(context.Context, ComponentConfig) error
	ShutdownComponent(context.Context, ComponentConfig) error
	KillComponent(context.Context, ComponentConfig) error
}

// XXX Method: Identify // XXX replace with side-effect of frefresh in the resource-component controller that returns a wrapper object with an iri field.
