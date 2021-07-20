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
}
