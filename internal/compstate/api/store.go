// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/internal/josh/server"
)

type Store interface {
	SetState(context.Context, *SetStateInput) (*SetStateOutput, error)
	GetStates(context.Context, *GetStatesInput) (*GetStatesOutput, error)
}

type SetStateInput struct {
	ComponentID string            `json:"componentId"`
	Type        string            `json:"type"`
	Content     string            `json:"content"`
	Tags        map[string]string `json:"tags"`
	Timestamp   string            `json:"timestamp"`
}

type SetStateOutput struct {
	Version int `json:"version"`
}

type GetStatesInput struct {
	ComponentID string `json:"componentId"`
	// If not specified, begins history with most recent.
	Version int `json:"version"`
	// Limit of historical states to return per component. Defaults to 1.
	History int `json:"history"`
}

type GetStatesOutput struct {

	// With descending version numbers.
	States []State `json:"states"`
}

func BuildStoreMux(b *josh.MuxBuilder, factory func(req *http.Request) Store) {
	b.AddMethod("set-state", func(req *http.Request) interface{} {
		return factory(req).SetState
	})
	b.AddMethod("get-states", func(req *http.Request) interface{} {
		return factory(req).GetStates
	})
}

type State struct {
	ComponentID string            `json:"componentId"`
	Version     int               `json:"version"`
	Type        string            `json:"type"`
	Content     string            `json:"content"`
	Tags        map[string]string `json:"tags"`
	Timestamp   string            `json:"timestamp"`
}
