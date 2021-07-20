// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

type Lifecycle interface {
	Initialize(context.Context, *InitializeInput) (*InitializeOutput, error)
	Update(context.Context, *UpdateInput) (*UpdateOutput, error)
	Refresh(context.Context, *RefreshInput) (*RefreshOutput, error)
	Dispose(context.Context, *DisposeInput) (*DisposeOutput, error)
}

type InitializeInput struct {
	ID   string `json:"id"`
	Spec string `json:"spec"`
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
	Spec  string `json:"spec"`
	State string `json:"state"`
}

type DisposeOutput struct {
	State string `json:"state"`
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
