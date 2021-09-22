// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/internal/josh/server"
)

type Kernel interface {
	CreateProject(context.Context, *CreateProjectInput) (*CreateProjectOutput, error)
	ListTemplates(context.Context, *ListTemplatesInput) (*ListTemplatesOutput, error)
	CreateWorkspace(context.Context, *CreateWorkspaceInput) (*CreateWorkspaceOutput, error)
	DescribeWorkspaces(context.Context, *DescribeWorkspacesInput) (*DescribeWorkspacesOutput, error)
	ResolveWorkspace(context.Context, *ResolveWorkspaceInput) (*ResolveWorkspaceOutput, error)
	// Debug method to test what happens when the service panics.
	Panic(context.Context, *PanicInput) (*PanicOutput, error)
	// Retrieves the installed and current version of exo.
	GetVersion(context.Context, *GetVersionInput) (*GetVersionOutput, error)
	// Upgrades exo to the latest version.
	Upgrade(context.Context, *UpgradeInput) (*UpgradeOutput, error)
	// Checks whether server is up.
	Ping(context.Context, *PingInput) (*PingOutput, error)
	// Gracefully shutdown the exo daemon.
	Exit(context.Context, *ExitInput) (*ExitOutput, error)
	DescribeTasks(context.Context, *DescribeTasksInput) (*DescribeTasksOutput, error)
}

type CreateProjectInput struct {
	Root        string  `json:"root"`
	TemplateUrl *string `json:"templateUrl"`
}

type CreateProjectOutput struct {
	WorkspaceID string `json:"workspaceId"`
}

type ListTemplatesInput struct {
}

type ListTemplatesOutput struct {
	TemplateNames []string `json:"templateNames"`
}

type CreateWorkspaceInput struct {
	Root string `json:"root"`
}

type CreateWorkspaceOutput struct {
	ID string `json:"id"`
}

type DescribeWorkspacesInput struct {
}

type DescribeWorkspacesOutput struct {
	Workspaces []WorkspaceDescription `json:"workspaces"`
}

type ResolveWorkspaceInput struct {
	Ref string `json:"ref"`
}

type ResolveWorkspaceOutput struct {
	ID *string `json:"id"`
}

type PanicInput struct {
	Message string `json:"message"`
}

type PanicOutput struct {
}

type GetVersionInput struct {
}

type GetVersionOutput struct {
	Installed string  `json:"installed"`
	Latest    *string `json:"latest"`
	Current   bool    `json:"current"`
}

type UpgradeInput struct {
}

type UpgradeOutput struct {
}

type PingInput struct {
}

type PingOutput struct {
}

type ExitInput struct {
}

type ExitOutput struct {
}

type DescribeTasksInput struct {

	// If supplied, filters tasks by job.
	JobIDs []string `json:"jobIds"`
}

type DescribeTasksOutput struct {
	Tasks []TaskDescription `json:"tasks"`
}

func BuildKernelMux(b *josh.MuxBuilder, factory func(req *http.Request) Kernel) {
	b.AddMethod("create-project", func(req *http.Request) interface{} {
		return factory(req).CreateProject
	})
	b.AddMethod("list-templates", func(req *http.Request) interface{} {
		return factory(req).ListTemplates
	})
	b.AddMethod("create-workspace", func(req *http.Request) interface{} {
		return factory(req).CreateWorkspace
	})
	b.AddMethod("describe-workspaces", func(req *http.Request) interface{} {
		return factory(req).DescribeWorkspaces
	})
	b.AddMethod("resolve-workspace", func(req *http.Request) interface{} {
		return factory(req).ResolveWorkspace
	})
	b.AddMethod("panic", func(req *http.Request) interface{} {
		return factory(req).Panic
	})
	b.AddMethod("get-version", func(req *http.Request) interface{} {
		return factory(req).GetVersion
	})
	b.AddMethod("upgrade", func(req *http.Request) interface{} {
		return factory(req).Upgrade
	})
	b.AddMethod("ping", func(req *http.Request) interface{} {
		return factory(req).Ping
	})
	b.AddMethod("exit", func(req *http.Request) interface{} {
		return factory(req).Exit
	})
	b.AddMethod("describe-tasks", func(req *http.Request) interface{} {
		return factory(req).DescribeTasks
	})
}

type TaskDescription struct {
	ID string `json:"id"`
	// ID of root task in this tree.
	JobID    string  `json:"jobId"`
	ParentID *string `json:"parentId"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	// Most recent log message. Single-line of text.
	Message  string        `json:"message"`
	Created  string        `json:"created"`
	Updated  string        `json:"updated"`
	Started  *string       `json:"started"`
	Finished *string       `json:"finished"`
	Progress *TaskProgress `json:"progress"`
}

type TaskProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
