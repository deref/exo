// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

type Kernel interface {
	CreateWorkspace(context.Context, *CreateWorkspaceInput) (*CreateWorkspaceOutput, error)
	DescribeWorkspaces(context.Context, *DescribeWorkspacesInput) (*DescribeWorkspacesOutput, error)
	FindWorkspace(context.Context, *FindWorkspaceInput) (*FindWorkspaceOutput, error)
	// Debug method to test what happens when the service panics.
	Panic(context.Context, *PanicInput) (*PanicOutput, error)
	// Retrieves the installed and current version of exo.
	GetVersion(context.Context, *GetVersionInput) (*GetVersionOutput, error)
	// Upgrades exo to the latest version.
	Upgrade(context.Context, *UpgradeInput) (*UpgradeOutput, error)
	// Checks whether server is up.
	Ping(context.Context, *PingInput) (*PingOutput, error)
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

type FindWorkspaceInput struct {
	Path string `json:"path"`
}

type FindWorkspaceOutput struct {
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

func BuildKernelMux(b *josh.MuxBuilder, factory func(req *http.Request) Kernel) {
	b.AddMethod("create-workspace", func(req *http.Request) interface{} {
		return factory(req).CreateWorkspace
	})
	b.AddMethod("describe-workspaces", func(req *http.Request) interface{} {
		return factory(req).DescribeWorkspaces
	})
	b.AddMethod("find-workspace", func(req *http.Request) interface{} {
		return factory(req).FindWorkspace
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
}
