// Generated file. DO NOT EDIT.

package api

import (
	"context"
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

	// ID of component to be patched.
	ID string `json:"id"`
	// If provided, renames component.
	Name string `json:"name"`
	// If provided, replaces component spec.
	Spec      string    `json:"spec"`
	State     string    `json:"state"`
	DependsOn *[]string `json:"dependsOn"`
}

type PatchComponentOutput struct {
}

type RemoveComponentInput struct {
	ID string `json:"id"`
}

type RemoveComponentOutput struct {
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
	DependsOn   []string `json:"dependsOn"`
}
