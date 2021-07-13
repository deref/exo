package api

import (
	"context"
	"net/http"

	"github.com/deref/exo/config"
	"github.com/deref/exo/josh"
)

type Project interface {
	// Deletes all of the components in the project, then deletes the project itself.
	Delete(context.Context, *DeleteInput) (*DeleteOutput, error)
	// Performs creates, updates, refreshes, disposes, as needed.
	Apply(context.Context, *ApplyInput) (*ApplyOutput, error)

	// Resolves a reference in to an ID.
	Resolve(context.Context, *ResolveInput) (*ResolveOutput, error)

	// Returns component descriptions.
	DescribeComponents(context.Context, *DescribeComponentsInput) (*DescribeComponentsOutput, error)
	// Creates a component and triggers an initialize lifecycle event.
	CreateComponent(context.Context, *CreateComponentInput) (*CreateComponentOutput, error)
	// Replaces the spec on a component and triggers an update lifecycle event.
	UpdateComponent(context.Context, *UpdateComponentInput) (*UpdateComponentOutput, error)
	// Triggers a refresh lifecycle event to update the component's state.
	RefreshComponent(context.Context, *RefreshComponentInput) (*RefreshComponentOutput, error)
	// Marks a component as disposed and triggers the dispose lifecycle event.
	// After being disposed, the component record will be deleted asynchronously.
	DisposeComponent(context.Context, *DisposeComponentInput) (*DisposeComponentOutput, error)
	// Disposes a component and then awaits the record to be deleted synchronously.
	DeleteComponent(context.Context, *DeleteComponentInput) (*DeleteComponentOutput, error)

	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
}

type DeleteInput struct{}

type DeleteOutput struct{}

type ApplyInput struct {
	Config config.Config `json:"config"`
}

type ApplyOutput struct{}

type ResolveInput struct {
	Refs []string `json:"refs"`
}

type ResolveOutput struct {
	IDs []*string `json:"ids"`
}

type DescribeComponentsInput struct{}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
}

type ComponentDescription struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Spec        string  `json:"spec"`
	State       string  `json:"state"`
	Created     string  `json:"created"`
	Initialized *string `json:"initialized"`
	Disposed    *string `json:"disposed"`
}

type CreateComponentInput struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Spec string `json:"spec"`
}

type CreateComponentOutput struct {
	ID string `json:"id"`
}

type UpdateComponentInput struct {
	Ref  string `json:"ref"`
	Spec string `json:"spec"`
}

type UpdateComponentOutput struct{}

type DisposeComponentInput struct {
	Ref string `json:"ref"`
}

type RefreshComponentInput struct {
	Ref string `json:"ref"`
}

type RefreshComponentOutput struct{}

type DisposeComponentOutput struct{}

type DeleteComponentInput struct {
	Ref string `json:"ref"`
}

type DeleteComponentOutput struct{}

type DescribeLogsInput struct {
	Refs []string `json:"refs"`
}

type DescribeLogsOutput struct {
	Logs []LogDescription `json:"logs"`
}

type LogDescription struct {
	Name        string  `json:"name"`
	LastEventAt *string `json:"lastEventAt"`
}

type GetEventsInput struct {
	Logs   []string `json:"logs"`
	Before string   `json:"before"`
	After  string   `json:"after"`
}

type GetEventsOutput struct {
	Events []Event `json:"events"`
}

type Event struct {
	Log       string `json:"log"`
	SID       string `json:"sid"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func NewProjectMux(prefix string, project Project) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(prefix+"delete", josh.NewMethodHandler(project.Delete))
	mux.Handle(prefix+"apply", josh.NewMethodHandler(project.Apply))
	mux.Handle(prefix+"resolve", josh.NewMethodHandler(project.Resolve))
	mux.Handle(prefix+"describe-components", josh.NewMethodHandler(project.DescribeComponents))
	mux.Handle(prefix+"create-component", josh.NewMethodHandler(project.CreateComponent))
	mux.Handle(prefix+"update-component", josh.NewMethodHandler(project.UpdateComponent))
	mux.Handle(prefix+"refresh-component", josh.NewMethodHandler(project.RefreshComponent))
	mux.Handle(prefix+"dispose-component", josh.NewMethodHandler(project.DisposeComponent))
	mux.Handle(prefix+"delete-component", josh.NewMethodHandler(project.DeleteComponent))
	mux.Handle(prefix+"describe-logs", josh.NewMethodHandler(project.DescribeLogs))
	mux.Handle(prefix+"get-events", josh.NewMethodHandler(project.GetEvents))
	return mux
}
