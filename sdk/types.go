package sdk

import (
	"context"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/api"
)

// Concrete component configuration.
// Corresponds to $Component in the schema.cue file.
type RawComponentConfig struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`

	// Intentionally not called "Spec" and "State". Those names are reserved for
	// more strongly typed field in Model structs that embed a ComponentConfig.
	SpecValue  cue.Value      `json:"spec"`
	StateValue map[string]any `json:"state"` // XXX reconsider how this works.

	Run         bool
	Environment map[string]string `json:"environment"`

	Resources map[string]ResourceConfig
}

type ResourceConfig struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	IRI  string `json:"iri,omitempty"`
}

type ComponentConfig[M Model] struct {
	RawComponentConfig
	Spec Model
}

type Controller[M Model] interface {
	InitializeController(context.Context, api.Service) error
	InitializeComponent(context.Context, ComponentConfig[M]) error
	ComponentCreated(context.Context, ComponentConfig[M]) error
	RenderComponent(context.Context, ComponentConfig[M]) ([]RenderedComponent, error)
	RefreshComponent(context.Context, ComponentConfig[M]) error
	ComponentUpdated(context.Context, ComponentConfig[M]) error
	ChildrenUpdated(context.Context, ComponentConfig[M]) error
	ShutdownComponent(context.Context, ComponentConfig[M]) error
	KillComponent(context.Context, ComponentConfig[M]) error
}

type Model interface {
	DecodeCue(v cue.Value) error
	EncodeCue() (string, error)
}

type RenderedComponent struct {
	Type string
	Name string
	Key  string
	Spec any
}
