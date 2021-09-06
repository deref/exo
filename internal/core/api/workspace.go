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
	Signal(context.Context, *SignalInput) (*SignalOutput, error)
	Restart(context.Context, *RestartInput) (*RestartOutput, error)
}

type StartInput struct {
}

type StartOutput struct {
	JobID string `json:"jobId"`
}

type StopInput struct {
	TimeoutSeconds *uint `json:"timeoutSeconds"`
}

type StopOutput struct {
	JobID string `json:"jobId"`
}

type SignalInput struct {
	Signal string `json:"signal"`
}

type SignalOutput struct {
	JobID string `json:"jobId"`
}

type RestartInput struct {
	TimeoutSeconds *uint `json:"timeoutSeconds"`
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
	b.AddMethod("signal", func(req *http.Request) interface{} {
		return factory(req).Signal
	})
	b.AddMethod("restart", func(req *http.Request) interface{} {
		return factory(req).Restart
	})
}

type Builder interface {
	Build(context.Context, *BuildInput) (*BuildOutput, error)
}

type BuildInput struct {
}

type BuildOutput struct {
	JobID string `json:"jobId"`
}

func BuildBuilderMux(b *josh.MuxBuilder, factory func(req *http.Request) Builder) {
	b.AddMethod("build", func(req *http.Request) interface{} {
		return factory(req).Build
	})
}

type Workspace interface {
	Process
	Builder
	// Describes this workspace.
	Describe(context.Context, *DescribeInput) (*DescribeOutput, error)
	// Asynchronously deletes all components in the workspace, then deletes the workspace itself.
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
	// Asynchronously runs dispose lifecycle methods on each component.
	DisposeComponents(context.Context, *DisposeComponentsInput) (*DisposeComponentsOutput, error)
	// Asynchronously disposes components, then removes them from the manifest.
	DeleteComponents(context.Context, *DeleteComponentsInput) (*DeleteComponentsOutput, error)
	GetComponentState(context.Context, *GetComponentStateInput) (*GetComponentStateOutput, error)
	SetComponentState(context.Context, *SetComponentStateInput) (*SetComponentStateOutput, error)
	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	// Returns pages of log events for some set of logs. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the log.
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
	StartComponents(context.Context, *StartComponentsInput) (*StartComponentsOutput, error)
	StopComponents(context.Context, *StopComponentsInput) (*StopComponentsOutput, error)
	SignalComponents(context.Context, *SignalComponentsInput) (*SignalComponentsOutput, error)
	RestartComponents(context.Context, *RestartComponentsInput) (*RestartComponentsOutput, error)
	DescribeProcesses(context.Context, *DescribeProcessesInput) (*DescribeProcessesOutput, error)
	DescribeVolumes(context.Context, *DescribeVolumesInput) (*DescribeVolumesOutput, error)
	DescribeNetworks(context.Context, *DescribeNetworksInput) (*DescribeNetworksOutput, error)
	ExportProcfile(context.Context, *ExportProcfileInput) (*ExportProcfileOutput, error)
	// Read a file from disk.
	ReadFile(context.Context, *ReadFileInput) (*ReadFileOutput, error)
	// Writes a file to disk.
	WriteFile(context.Context, *WriteFileInput) (*WriteFileOutput, error)
	BuildComponents(context.Context, *BuildComponentsInput) (*BuildComponentsOutput, error)
	DescribeEnvironment(context.Context, *DescribeEnvironmentInput) (*DescribeEnvironmentOutput, error)
}

type DescribeInput struct {
}

type DescribeOutput struct {
	Description WorkspaceDescription `json:"description"`
}

type DestroyInput struct {
}

type DestroyOutput struct {
	JobID string `json:"jobId"`
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
	JobID    string   `json:"jobId"`
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
	// If true, includes all components that the filtered components depend on.
	IncludeDependencies bool `json:"includeDependencies"`
	// If true, includes all components that depend on the filtered components.
	IncludeDependents bool `json:"includeDependents"`
}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
}

type CreateComponentInput struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Spec      string   `json:"spec"`
	DependsOn []string `json:"dependsOn"`
}

type CreateComponentOutput struct {
	ID string `json:"id"`
}

type UpdateComponentInput struct {
	Ref       string   `json:"ref"`
	Spec      string   `json:"spec"`
	DependsOn []string `json:"dependsOn"`
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

type DisposeComponentsInput struct {
	Refs []string `json:"refs"`
}

type DisposeComponentsOutput struct {
	JobID string `json:"jobId"`
}

type DeleteComponentsInput struct {
	Refs []string `json:"refs"`
}

type DeleteComponentsOutput struct {
	JobID string `json:"jobId"`
}

type GetComponentStateInput struct {
	Ref string `json:"ref"`
}

type GetComponentStateOutput struct {
	State string `json:"state"`
}

type SetComponentStateInput struct {
	Ref   string `json:"ref"`
	State string `json:"state"`
}

type SetComponentStateOutput struct {
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
	Refs           []string `json:"refs"`
	TimeoutSeconds *uint    `json:"timeoutSeconds"`
}

type StopComponentsOutput struct {
	JobID string `json:"jobId"`
}

type SignalComponentsInput struct {
	Refs   []string `json:"refs"`
	Signal string   `json:"signal"`
}

type SignalComponentsOutput struct {
	JobID string `json:"jobId"`
}

type RestartComponentsInput struct {
	Refs           []string `json:"refs"`
	TimeoutSeconds *uint    `json:"timeoutSeconds"`
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

type ExportProcfileInput struct {
}

type ExportProcfileOutput struct {
	Procfile string `json:"procfile"`
}

type ReadFileInput struct {

	// Relative to the workspace directory. May not traverse higher in the filesystem.
	Path string `json:"path"`
}

type ReadFileOutput struct {
	Content string `json:"content"`
}

type WriteFileInput struct {

	// Relative to the workspace directory. May not traverse higher in the filesystem.
	Path    string `json:"path"`
	Mode    *int   `json:"mode"`
	Content string `json:"content"`
}

type WriteFileOutput struct {
}

type BuildComponentsInput struct {
	Refs []string `json:"refs"`
}

type BuildComponentsOutput struct {
	JobID string `json:"jobId"`
}

type DescribeEnvironmentInput struct {
}

type DescribeEnvironmentOutput struct {
	Variables map[string]string `json:"variables"`
}

func BuildWorkspaceMux(b *josh.MuxBuilder, factory func(req *http.Request) Workspace) {
	b.AddMethod("start", func(req *http.Request) interface{} {
		return factory(req).Start
	})
	b.AddMethod("stop", func(req *http.Request) interface{} {
		return factory(req).Stop
	})
	b.AddMethod("signal", func(req *http.Request) interface{} {
		return factory(req).Signal
	})
	b.AddMethod("restart", func(req *http.Request) interface{} {
		return factory(req).Restart
	})
	b.AddMethod("build", func(req *http.Request) interface{} {
		return factory(req).Build
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
	b.AddMethod("dispose-components", func(req *http.Request) interface{} {
		return factory(req).DisposeComponents
	})
	b.AddMethod("delete-components", func(req *http.Request) interface{} {
		return factory(req).DeleteComponents
	})
	b.AddMethod("get-component-state", func(req *http.Request) interface{} {
		return factory(req).GetComponentState
	})
	b.AddMethod("set-component-state", func(req *http.Request) interface{} {
		return factory(req).SetComponentState
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
	b.AddMethod("signal-components", func(req *http.Request) interface{} {
		return factory(req).SignalComponents
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
	b.AddMethod("export-procfile", func(req *http.Request) interface{} {
		return factory(req).ExportProcfile
	})
	b.AddMethod("read-file", func(req *http.Request) interface{} {
		return factory(req).ReadFile
	})
	b.AddMethod("write-file", func(req *http.Request) interface{} {
		return factory(req).WriteFile
	})
	b.AddMethod("build-components", func(req *http.Request) interface{} {
		return factory(req).BuildComponents
	})
	b.AddMethod("describe-environment", func(req *http.Request) interface{} {
		return factory(req).DescribeEnvironment
	})
}

type WorkspaceDescription struct {
	ID   string `json:"id"`
	Root string `json:"root"`
}

type ComponentDescription struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Spec        string   `json:"spec"`
	State       string   `json:"state"`
	Created     string   `json:"created"`
	Initialized *string  `json:"initialized"`
	Disposed    *string  `json:"disposed"`
	DependsOn   []string `json:"dependsOn"`
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
	Spec                string            `json:"spec"`
	Running             bool              `json:"running"`
	EnvVars             map[string]string `json:"envVars"`
	CPUPercent          *float64          `json:"cpuPercent"`
	CreateTime          *int64            `json:"createTime"`
	ResidentMemory      *uint64           `json:"residentMemory"`
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
