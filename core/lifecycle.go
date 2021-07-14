// TODO: Generate this package via JOSH introspection.

package core

import (
	"context"
)

type Lifecycle interface {
	Initialize(context.Context, *InitializeInput) (*InitializeOutput, error)
	Update(context.Context, *UpdateInput) (*UpdateOutput, error)
	Refresh(context.Context, *RefreshInput) (*RefreshOutput, error)
	Dispose(context.Context, *DisposeInput) (*DisposeOutput, error)
}

type InitializeInput struct {
	ID   string `json:"id"`
	Spec string `json:"spec"` // TODO: content-type tagged data, default to application/json or whatever.
}

type InitializeOutput struct {
	State string `json:"state"`
}

type UpdateInput struct {
	ID      string `json:"id"`
	OldSpec string `json:"oldSpec"`
	NewSpec string `json:"newSpec"`
	State   string `json:"state"`
}

type UpdateOutput struct {
	State string `json:"state"`
}

type RefreshInput struct {
	ID    string `json:"id"`
	Spec  string `json:"spec"`
	State string `json:"state"`
}

type RefreshOutput struct {
	State string `json:"state"`
}

type DisposeInput struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

type DisposeOutput struct {
	State string `json:"state"`
	// TODO: Return a promise that can be awaited for synchronous deletes.
}
