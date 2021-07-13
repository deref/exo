// TODO: Generate this package via JOSH introspection.

package api

import (
	"context"
	"net/http"

	"github.com/deref/exo/josh"
)

type Lifecycle interface {
	Initialize(context.Context, *InitializeInput) (*InitializeOutput, error)
	Update(context.Context, *UpdateInput) (*UpdateOutput, error)
	Refresh(context.Context, *RefreshInput) (*RefreshOutput, error)
	Dispose(context.Context, *DisposeInput) (*DisposeOutput, error)
}

type InitializeInput struct {
	ID   string                 `json:"id"`
	Spec map[string]interface{} `json:"spec"` // TODO: content-type tagged data, default to application/json or whatever.
}

type InitializeOutput struct {
	State map[string]interface{} `json:"state"`
}

type UpdateInput struct {
	ID      string                 `json:"id"`
	OldSpec map[string]interface{} `json:"oldSpec"`
	NewSpec map[string]interface{} `json:"newSpec"`
	State   map[string]interface{} `json:"state"`
}

type UpdateOutput struct {
	State map[string]interface{} `json:"state"`
}

type RefreshInput struct {
	ID    string                 `json:"id"`
	State map[string]interface{} `json:"state"`
}

type RefreshOutput struct {
	State map[string]interface{} `json:"state"`
}

type DisposeInput struct {
	ID    string                 `json:"id"`
	State map[string]interface{} `json:"state"`
}

type DisposeOutput struct {
	State map[string]interface{} `json:"state"`
}

func NewLifecycleMux(prefix string, lifecycle Lifecycle) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(prefix+"initialize", josh.NewMethodHandler(lifecycle.Initialize))
	mux.Handle(prefix+"update", josh.NewMethodHandler(lifecycle.Update))
	mux.Handle(prefix+"refresh", josh.NewMethodHandler(lifecycle.Refresh))
	mux.Handle(prefix+"dispose", josh.NewMethodHandler(lifecycle.Dispose))
	return mux
}
