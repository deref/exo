// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/internal/josh/server"
)

type Store interface {
	// Returns workspace descriptions.
	DescribeWorkspaces(context.Context, *DescribeWorkspacesInput) (*DescribeWorkspacesOutput, error)
	AddWorkspace(context.Context, *AddWorkspaceInput) (*AddWorkspaceOutput, error)
	RemoveWorkspace(context.Context, *RemoveWorkspaceInput) (*RemoveWorkspaceOutput, error)
	ResolveWorkspace(context.Context, *ResolveWorkspaceInput) (*ResolveWorkspaceOutput, error)
	Resolve(context.Context, *ResolveInput) (*ResolveOutput, error)
	DescribeComponents(context.Context, *DescribeComponentsInput) (*DescribeComponentsOutput, error)
	AddComponent(context.Context, *AddComponentInput) (*AddComponentOutput, error)
	PatchComponent(context.Context, *PatchComponentInput) (*PatchComponentOutput, error)
	RemoveComponent(context.Context, *RemoveComponentInput) (*RemoveComponentOutput, error)
}

type DescribeWorkspacesInput struct {
	IDs []string `json:"ids"`
}

type DescribeWorkspacesOutput struct {
	Workspaces []WorkspaceDescription `json:"workspaces"`
}

type AddWorkspaceInput struct {
	ID   string `json:"id"`
	Root string `json:"root"`
}

type AddWorkspaceOutput struct {
}

type RemoveWorkspaceInput struct {
	ID string `json:"id"`
}

type RemoveWorkspaceOutput struct {
}

type ResolveWorkspaceInput struct {
	Ref string `json:"ref"`
}

type ResolveWorkspaceOutput struct {
	ID *string `json:"id"`
}

type ResolveInput struct {
	WorkspaceID string   `json:"workspaceId"`
	Refs        []string `json:"refs"`
}

type ResolveOutput struct {
	IDs []*string `json:"ids"`
}

type DescribeComponentsInput struct {
	WorkspaceID         string   `json:"workspaceId"`
	Refs                []string `json:"refs"`
	Types               []string `json:"types"`
	IncludeDependencies bool     `json:"includeDependencies"`
	IncludeDependents   bool     `json:"includeDependents"`
}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
}

type AddComponentInput struct {
	WorkspaceID string   `json:"workspaceId"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Spec        string   `json:"spec"`
	Created     string   `json:"created"`
	DependsOn   []string `json:"dependsOn"`
}

type AddComponentOutput struct {
}

type PatchComponentInput struct {
	ID          string    `json:"id"`
	State       string    `json:"state"`
	Initialized string    `json:"initialized"`
	Disposed    string    `json:"disposed"`
	DependsOn   *[]string `json:"dependsOn"`
}

type PatchComponentOutput struct {
}

type RemoveComponentInput struct {
	ID string `json:"id"`
}

type RemoveComponentOutput struct {
}

func BuildStoreMux(b *josh.MuxBuilder, factory func(req *http.Request) Store) {
	b.AddMethod("describe-workspaces", func(req *http.Request) interface{} {
		return factory(req).DescribeWorkspaces
	})
	b.AddMethod("add-workspace", func(req *http.Request) interface{} {
		return factory(req).AddWorkspace
	})
	b.AddMethod("remove-workspace", func(req *http.Request) interface{} {
		return factory(req).RemoveWorkspace
	})
	b.AddMethod("resolve-workspace", func(req *http.Request) interface{} {
		return factory(req).ResolveWorkspace
	})
	b.AddMethod("resolve", func(req *http.Request) interface{} {
		return factory(req).Resolve
	})
	b.AddMethod("describe-components", func(req *http.Request) interface{} {
		return factory(req).DescribeComponents
	})
	b.AddMethod("add-component", func(req *http.Request) interface{} {
		return factory(req).AddComponent
	})
	b.AddMethod("patch-component", func(req *http.Request) interface{} {
		return factory(req).PatchComponent
	})
	b.AddMethod("remove-component", func(req *http.Request) interface{} {
		return factory(req).RemoveComponent
	})
}

type WorkspaceDescription struct {
	ID          string `json:"id"`
	Root        string `json:"root"`
	DisplayName string `json:"displayName"`
}

type ComponentDescription struct {
	ID          string   `json:"id"`
	WorkspaceID string   `json:"workspaceId"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Spec        string   `json:"spec"`
	State       string   `json:"state"`
	Created     string   `json:"created"`
	Initialized *string  `json:"initialized"`
	Disposed    *string  `json:"disposed"`
	DependsOn   []string `json:"dependsOn"`
}
