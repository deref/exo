// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

type Project interface {
	// Deletes all of the components in the project, then deletes the project itself.
	Delete(context.Context, *DeleteInput) (*DeleteOutput, error)
	// Performs creates, updates, refreshes, disposes, as needed.
	Apply(context.Context, *ApplyInput) (*ApplyOutput, error)
	// Refreshes all components.
	Refresh(context.Context, *RefreshInput) (*RefreshOutput, error)
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
	// Marks a component as disposed and triggers the dispose lifecycle event. After being disposed, the component record will be deleted asynchronously.
	DisposeComponent(context.Context, *DisposeComponentInput) (*DisposeComponentOutput, error)
	// Disposes a component and then awaits the record to be deleted synchronously.
	DeleteComponent(context.Context, *DeleteComponentInput) (*DeleteComponentOutput, error)
	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
	Start(context.Context, *StartInput) (*StartOutput, error)
	Stop(context.Context, *StopInput) (*StopOutput, error)
	DescribeProcesses(context.Context, *DescribeProcessesInput) (*DescribeProcessesOutput, error)
}

type DeleteInput struct {
}

type DeleteOutput struct {
}

type ApplyInput struct {
	Config string `json:"config"`
}

type ApplyOutput struct {
}

type RefreshInput struct {
}

type RefreshOutput struct {
}

type ResolveInput struct {
	Refs []string `json:"refs"`
}

type ResolveOutput struct {
	IDs []*string `json:"ids"`
}

type DescribeComponentsInput struct {
}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
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

type UpdateComponentOutput struct {
}

type RefreshComponentInput struct {
	Ref string `json:"ref"`
}

type RefreshComponentOutput struct {
}

type DisposeComponentInput struct {
	Ref string `json:"ref"`
}

type DisposeComponentOutput struct {
}

type DeleteComponentInput struct {
	Ref string `json:"ref"`
}

type DeleteComponentOutput struct {
}

type DescribeLogsInput struct {
	Refs []string `json:"refs"`
}

type DescribeLogsOutput struct {
	Logs []LogDescription `json:"logs"`
}

type GetEventsInput struct {
	Logs   []string `json:"logs"`
	Before string   `json:"before"`
	After  string   `json:"after"`
}

type GetEventsOutput struct {
	Events []Event `json:"events"`
}

type StartInput struct {
	Ref string `json:"ref"`
}

type StartOutput struct {
}

type StopInput struct {
	Ref string `json:"ref"`
}

type StopOutput struct {
}

type DescribeProcessesInput struct {
}

type DescribeProcessesOutput struct {
	Processes []ProcessDescription `json:"processes"`
}

func NewProjectMux(prefix string, iface Project) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	BuildProjectMux(b, iface)
	return b.Mux()
}

func BuildProjectMux(b *josh.MuxBuilder, iface Project) {
	b.AddMethod("delete", iface.Delete)
	b.AddMethod("apply", iface.Apply)
	b.AddMethod("refresh", iface.Refresh)
	b.AddMethod("resolve", iface.Resolve)
	b.AddMethod("describe-components", iface.DescribeComponents)
	b.AddMethod("create-component", iface.CreateComponent)
	b.AddMethod("update-component", iface.UpdateComponent)
	b.AddMethod("refresh-component", iface.RefreshComponent)
	b.AddMethod("dispose-component", iface.DisposeComponent)
	b.AddMethod("delete-component", iface.DeleteComponent)
	b.AddMethod("describe-logs", iface.DescribeLogs)
	b.AddMethod("get-events", iface.GetEvents)
	b.AddMethod("start", iface.Start)
	b.AddMethod("stop", iface.Stop)
	b.AddMethod("describe-processes", iface.DescribeProcesses)
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

type LogDescription struct {
	Name        string  `json:"name"`
	LastEventAt *string `json:"lastEventAt"`
}

type Event struct {
	ID        string `json:"id"`
	Log       string `json:"log"`
	SID       string `json:"sid"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type ProcessDescription struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Running bool   `json:"running"`
}
