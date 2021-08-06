// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/internal/josh/server"
)

type Process interface {
	Start(context.Context, *StartInput) (*StartOutput, error)
	Stop(context.Context, *StopInput) (*StopOutput, error)
	Restart(context.Context, *RestartInput) (*RestartOutput, error)
}

type StartInput struct {
}

type StartOutput struct {
	JobID string `json:"jobId"`
}

type StopInput struct {
}

type StopOutput struct {
	JobID string `json:"jobId"`
}

type RestartInput struct {
}

type RestartOutput struct {
	JobID string `json:"jobId"`
}

func BuildProcessMux(b *josh.MuxBuilder, factory func(req *http.Request) Process) {
	b.AddMethod("start", func(req *http.Request) interface{} {
		return factory(req).Start
	})
	b.AddMethod("stop", func(req *http.Request) interface{} {
		return factory(req).Stop
	})
	b.AddMethod("restart", func(req *http.Request) interface{} {
		return factory(req).Restart
	})
}

type Workspace interface {
	Process
	// Describes this workspace.
	Describe(context.Context, *DescribeInput) (*DescribeOutput, error)
	// Deletes all of the components in the workspace, then deletes the workspace itself.
	Destroy(context.Context, *DestroyInput) (*DestroyOutput, error)
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
	// Asycnhronously refreshes component state.
	RefreshComponents(context.Context, *RefreshComponentsInput) (*RefreshComponentsOutput, error)
	// Marks a component as disposed and triggers the dispose lifecycle event. After being disposed, the component record will be deleted asynchronously.
	DisposeComponent(context.Context, *DisposeComponentInput) (*DisposeComponentOutput, error)
	// Disposes a component and then awaits the record to be deleted synchronously.
	DeleteComponent(context.Context, *DeleteComponentInput) (*DeleteComponentOutput, error)
	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	// Returns pages of log events for some set of logs. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the log.
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
	StartComponents(context.Context, *StartComponentsInput) (*StartComponentsOutput, error)
	StopComponents(context.Context, *StopComponentsInput) (*StopComponentsOutput, error)
	RestartComponents(context.Context, *RestartComponentsInput) (*RestartComponentsOutput, error)
	DescribeProcesses(context.Context, *DescribeProcessesInput) (*DescribeProcessesOutput, error)
	DescribeVolumes(context.Context, *DescribeVolumesInput) (*DescribeVolumesOutput, error)
	DescribeNetworks(context.Context, *DescribeNetworksInput) (*DescribeNetworksOutput, error)
}

type DescribeInput struct {
}

type DescribeOutput struct {
	Description WorkspaceDescription `json:"description"`
}

type DestroyInput struct {
}

type DestroyOutput struct {
}

type ApplyInput struct {

	// One of 'exo', 'compose', or 'procfile'.
	Format *string `json:"format"`
	// Path of manifest file to load. May be relative to the workspace root. If format is not provided, will be inferred from path name.
	ManifestPath *string `json:"manifestPath"`
	// Contents of the manifest file. Not required if manifest-path is provided.
	Manifest *string `json:"manifest"`
}

type ApplyOutput struct {
	Warnings []string `json:"warnings"`
}

type ResolveInput struct {
	Refs []string `json:"refs"`
}

type ResolveOutput struct {
	IDs []*string `json:"ids"`
}

type DescribeComponentsInput struct {

	// If non-empty, filters components to supplied ids.
	IDs []string `json:"ids"`
	// If non-empty, filters components to supplied types.
	Types []string `json:"types"`
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

type RefreshComponentsInput struct {

	// If omitted, refreshes all components.
	Refs []string `json:"refs"`
}

type RefreshComponentsOutput struct {
	JobID string `json:"jobId"`
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
	Logs      []string `json:"logs"`
	Cursor    *string  `json:"cursor"`
	FilterStr *string  `json:"filterStr"`
	Prev      *int     `json:"prev"`
	Next      *int     `json:"next"`
}

type GetEventsOutput struct {
	Items      []Event `json:"items"`
	PrevCursor string  `json:"prevCursor"`
	NextCursor string  `json:"nextCursor"`
}

type StartComponentsInput struct {
	Refs []string `json:"refs"`
}

type StartComponentsOutput struct {
	JobID string `json:"jobId"`
}

type StopComponentsInput struct {
	Refs []string `json:"refs"`
}

type StopComponentsOutput struct {
	JobID string `json:"jobId"`
}

type RestartComponentsInput struct {
	Refs []string `json:"refs"`
}

type RestartComponentsOutput struct {
	JobID string `json:"jobId"`
}

type DescribeProcessesInput struct {
}

type DescribeProcessesOutput struct {
	Processes []ProcessDescription `json:"processes"`
}

type DescribeVolumesInput struct {
}

type DescribeVolumesOutput struct {
	Volumes []VolumeDescription `json:"volumes"`
}

type DescribeNetworksInput struct {
}

type DescribeNetworksOutput struct {
	Networks []NetworkDescription `json:"networks"`
}

func BuildWorkspaceMux(b *josh.MuxBuilder, factory func(req *http.Request) Workspace) {
	b.AddMethod("start", func(req *http.Request) interface{} {
		return factory(req).Start
	})
	b.AddMethod("stop", func(req *http.Request) interface{} {
		return factory(req).Stop
	})
	b.AddMethod("restart", func(req *http.Request) interface{} {
		return factory(req).Restart
	})
	b.AddMethod("describe", func(req *http.Request) interface{} {
		return factory(req).Describe
	})
	b.AddMethod("destroy", func(req *http.Request) interface{} {
		return factory(req).Destroy
	})
	b.AddMethod("apply", func(req *http.Request) interface{} {
		return factory(req).Apply
	})
	b.AddMethod("resolve", func(req *http.Request) interface{} {
		return factory(req).Resolve
	})
	b.AddMethod("describe-components", func(req *http.Request) interface{} {
		return factory(req).DescribeComponents
	})
	b.AddMethod("create-component", func(req *http.Request) interface{} {
		return factory(req).CreateComponent
	})
	b.AddMethod("update-component", func(req *http.Request) interface{} {
		return factory(req).UpdateComponent
	})
	b.AddMethod("refresh-components", func(req *http.Request) interface{} {
		return factory(req).RefreshComponents
	})
	b.AddMethod("dispose-component", func(req *http.Request) interface{} {
		return factory(req).DisposeComponent
	})
	b.AddMethod("delete-component", func(req *http.Request) interface{} {
		return factory(req).DeleteComponent
	})
	b.AddMethod("describe-logs", func(req *http.Request) interface{} {
		return factory(req).DescribeLogs
	})
	b.AddMethod("get-events", func(req *http.Request) interface{} {
		return factory(req).GetEvents
	})
	b.AddMethod("start-components", func(req *http.Request) interface{} {
		return factory(req).StartComponents
	})
	b.AddMethod("stop-components", func(req *http.Request) interface{} {
		return factory(req).StopComponents
	})
	b.AddMethod("restart-components", func(req *http.Request) interface{} {
		return factory(req).RestartComponents
	})
	b.AddMethod("describe-processes", func(req *http.Request) interface{} {
		return factory(req).DescribeProcesses
	})
	b.AddMethod("describe-volumes", func(req *http.Request) interface{} {
		return factory(req).DescribeVolumes
	})
	b.AddMethod("describe-networks", func(req *http.Request) interface{} {
		return factory(req).DescribeNetworks
	})
}

type WorkspaceDescription struct {
	ID   string `json:"id"`
	Root string `json:"root"`
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
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type ProcessDescription struct {
	ID                  string            `json:"id"`
	Provider            string            `json:"provider"`
	Name                string            `json:"name"`
	Running             bool              `json:"running"`
	EnvVars             map[string]string `json:"envVars"`
	CPUPercent          float64           `json:"cpuPercent"`
	CreateTime          int64             `json:"createTime"`
	ResidentMemory      uint64            `json:"residentMemory"`
	Ports               []uint32          `json:"ports"`
	ChildrenExecutables []string          `json:"childrenExecutables"`
}

type VolumeDescription struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NetworkDescription struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
