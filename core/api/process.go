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

func NewProcessMux(prefix string, iface Process) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	BuildProcessMux(b, iface)
	return b.Mux()
}

func BuildProcessMux(b *josh.MuxBuilder, iface Process) {
	b.AddMethod("start", iface.Start)
	b.AddMethod("stop", iface.Stop)
}
