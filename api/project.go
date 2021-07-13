package api

import (
	"context"
	"net/http"

	"github.com/deref/exo/config"
	"github.com/deref/exo/josh"
)

type Project interface {
	// Returns component descriptions.
	DescribeComponents(context.Context, *DescribeComponentsInput) (*DescribeComponentsOutput, error)
	// Performs creates, updates, refreshes, disposes, as needed.
	Apply(context.Context, *ApplyInput) (*ApplyOutput, error)
	// Creates a component and triggers an initialize lifecycle event.
	CreateComponent(context.Context, *CreateComponentInput) (*CreateComponentOutput, error)
	// Replaces the spec on a component and triggers an update lifecycle event.
	UpdateComponent(context.Context, *UpdateComponentInput) (*UpdateComponentOutput, error)
	// Triggers a refresh lifecycle event to update the component's state.
	RefreshComponent(context.Context, *RefreshComponentInput) (*RefreshComponentOutput, error)
	// Marks a component as disposed and triggers the dispose lifecycle event.
	// Afterwards, the component will be deleted.
	DisposeComponent(context.Context, *DisposeComponentInput) (*DisposeComponentOutput, error)
	// Fails if not disposed unless you specify force.
	DeleteComponent(context.Context, *DeleteComponentInput) (*DeleteComponentOutput, error)
}

type DescribeComponentsInput struct{}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
}

type ComponentDescription struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Spec        map[string]interface{} `json:"spec"`
	State       map[string]interface{} `json:"state"`
	Created     string                 `json:"created"`
	Initialized *string                `json:"initialized"`
	Disposed    *string                `json:"disposed"`
}

type ApplyInput struct {
	Config config.Config `json:"config"`
}

type ApplyOutput struct{}

type CreateComponentInput struct {
	Name string                 `json:"name"`
	Type string                 `json:"type"`
	Spec map[string]interface{} `json:"spec"` // TODO: content-type tagged data, default to application/json or whatever.
}

type CreateComponentOutput struct {
	ID string `json:"id"`
}

type UpdateComponentInput struct {
	Name string                 `json:"name"`
	Spec map[string]interface{} `json:"spec"`
}

type UpdateComponentOutput struct{}

type DisposeComponentInput struct {
	Name string `json:"name"`
}

type RefreshComponentInput struct {
	Name string `json:"name"`
}

type RefreshComponentOutput struct{}

type DisposeComponentOutput struct{}

type DeleteComponentInput struct {
	Name  string `json:"name"`
	Force bool   `json:"force"`
}

type DeleteComponentOutput struct{}

func NewProjectMux(prefix string, project Project) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(prefix+"describe-components", josh.NewMethodHandler(project.DescribeComponents))
	mux.Handle(prefix+"apply", josh.NewMethodHandler(project.Apply))
	mux.Handle(prefix+"create-component", josh.NewMethodHandler(project.CreateComponent))
	mux.Handle(prefix+"update-component", josh.NewMethodHandler(project.UpdateComponent))
	mux.Handle(prefix+"refresh-component", josh.NewMethodHandler(project.RefreshComponent))
	mux.Handle(prefix+"dispose-component", josh.NewMethodHandler(project.DisposeComponent))
	mux.Handle(prefix+"delete-component", josh.NewMethodHandler(project.DeleteComponent))
	return mux
}
