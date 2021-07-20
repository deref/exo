// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

type Process interface {
	Start(context.Context, *StartInput) (*StartOutput, error)
	Stop(context.Context, *StopInput) (*StopOutput, error)
}

type StartInput struct {
	ID    string `json:"id"`
	Spec  string `json:"spec"`
	State string `json:"state"`
}

type StartOutput struct {
	State string `json:"state"`
}

type StopInput struct {
	ID    string `json:"id"`
	Spec  string `json:"spec"`
	State string `json:"state"`
}

type StopOutput struct {
	State string `json:"state"`
}

func BuildProcessMux(b *josh.MuxBuilder, factory func(req *http.Request) Process) {
	b.AddMethod("start", func(req *http.Request) interface{} {
		return factory(req).Start
	})
	b.AddMethod("stop", func(req *http.Request) interface{} {
		return factory(req).Stop
	})
}
