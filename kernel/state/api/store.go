// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

type Store interface {
	Resolve(context.Context, *ResolveInput) (*ResolveOutput, error)
	DescribeComponents(context.Context, *DescribeComponentsInput) (*DescribeComponentsOutput, error)
	AddComponent(context.Context, *AddComponentInput) (*AddComponentOutput, error)
	PatchComponent(context.Context, *PatchComponentInput) (*PatchComponentOutput, error)
	RemoveComponent(context.Context, *RemoveComponentInput) (*RemoveComponentOutput, error)
}

type ResolveInput struct {
	WorkspaceID string   `json:"workspaceId"`
	Refs        []string `json:"refs"`
}

type ResolveOutput struct {
	IDs []*string `json:"ids"`
}

type DescribeComponentsInput struct {
	WorkspaceID string   `json:"workspaceId"`
	IDs         []string `json:"ids"`
}

type DescribeComponentsOutput struct {
	Components []ComponentDescription `json:"components"`
}

type AddComponentInput struct {
	WorkspaceID string `json:"workspaceId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Spec        string `json:"spec"`
	Created     string `json:"created"`
}

type AddComponentOutput struct {
}

type PatchComponentInput struct {
	ID          string `json:"id"`
	State       string `json:"state"`
	Initialized string `json:"initialized"`
	Disposed    string `json:"disposed"`
}

type PatchComponentOutput struct {
}

type RemoveComponentInput struct {
	ID string `json:"id"`
}

type RemoveComponentOutput struct {
}

func NewStoreMux(prefix string, iface Store) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	BuildStoreMux(b, iface)
	return b.Mux()
}

func BuildStoreMux(b *josh.MuxBuilder, iface Store) {
	b.AddMethod("resolve", iface.Resolve)
	b.AddMethod("describe-components", iface.DescribeComponents)
	b.AddMethod("add-component", iface.AddComponent)
	b.AddMethod("patch-component", iface.PatchComponent)
	b.AddMethod("remove-component", iface.RemoveComponent)
}

type ComponentDescription struct {
	ID          string  `json:"id"`
	WorkspaceID string  `json:"workspaceId"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Spec        string  `json:"spec"`
	State       string  `json:"state"`
	Created     string  `json:"created"`
	Initialized *string `json:"initialized"`
	Disposed    *string `json:"disposed"`
}
