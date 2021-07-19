// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

type Kernel interface {
	CreateWorkspace(context.Context, *CreateWorkspaceInput) (*CreateWorkspaceOutput, error)
	ForgetWorkspace(context.Context, *ForgetWorkspaceInput) (*ForgetWorkspaceOutput, error)
	DescribeWorkspaces(context.Context, *DescribeWorkspacesInput) (*DescribeWorkspacesOutput, error)
	FindWorkspace(context.Context, *FindWorkspaceInput) (*FindWorkspaceOutput, error)
}

type CreateWorkspaceInput struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type CreateWorkspaceOutput struct {
}

type ForgetWorkspaceInput struct {
	Ref string `json:"ref"`
}

type ForgetWorkspaceOutput struct {
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

func NewKernelMux(prefix string, iface Kernel) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	BuildKernelMux(b, iface)
	return b.Mux()
}

func BuildKernelMux(b *josh.MuxBuilder, iface Kernel) {
	b.AddMethod("create-workspace", iface.CreateWorkspace)
	b.AddMethod("forget-workspace", iface.ForgetWorkspace)
	b.AddMethod("describe-workspaces", iface.DescribeWorkspaces)
	b.AddMethod("find-workspace", iface.FindWorkspace)
}

type WorkspaceDescription struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Root    string `json:"root"`
	Project string `json:"project"`
}
