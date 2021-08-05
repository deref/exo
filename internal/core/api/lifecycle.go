// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/internal/josh/server"
)

type Lifecycle interface {
	Initialize(context.Context, *InitializeInput) (*InitializeOutput, error)
	Update(context.Context, *UpdateInput) (*UpdateOutput, error)
	Refresh(context.Context, *RefreshInput) (*RefreshOutput, error)
	Dispose(context.Context, *DisposeInput) (*DisposeOutput, error)
}

type InitializeInput struct {
}

type InitializeOutput struct {
}

type UpdateInput struct {
	NewSpec string `json:"newSpec"`
}

type UpdateOutput struct {
}

type RefreshInput struct {
}

type RefreshOutput struct {
}

type DisposeInput struct {
}

type DisposeOutput struct {
}

func BuildLifecycleMux(b *josh.MuxBuilder, factory func(req *http.Request) Lifecycle) {
	b.AddMethod("initialize", func(req *http.Request) interface{} {
		return factory(req).Initialize
	})
	b.AddMethod("update", func(req *http.Request) interface{} {
		return factory(req).Update
	})
	b.AddMethod("refresh", func(req *http.Request) interface{} {
		return factory(req).Refresh
	})
	b.AddMethod("dispose", func(req *http.Request) interface{} {
		return factory(req).Dispose
	})
}
