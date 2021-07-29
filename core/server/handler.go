package server

import (
	"net/http"

	"github.com/deref/exo/core/api"
	state "github.com/deref/exo/core/state/api"
	josh "github.com/deref/exo/josh/server"
	"github.com/deref/exo/telemetry"
)

type Config struct {
	VarDir     string
	Store      state.Store
	Telemetry  telemetry.Telemetry
	SyslogAddr string
}

func BuildRootMux(prefix string, cfg *Config) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)

	endKernel := b.Begin("kernel")
	api.BuildKernelMux(b, func(req *http.Request) api.Kernel {
		return &Kernel{
			VarDir:    cfg.VarDir,
			Store:     cfg.Store,
			Telemetry: cfg.Telemetry,
		}
	})
	endKernel()

	endWorkspace := b.Begin("workspace")
	api.BuildWorkspaceMux(b, func(req *http.Request) api.Workspace {
		return &Workspace{
			ID:         req.URL.Query().Get("id"),
			VarDir:     cfg.VarDir,
			Store:      cfg.Store,
			SyslogAddr: cfg.SyslogAddr,
		}
	})
	endWorkspace()

	mux := b.Build()

	mux.Handle(prefix+"health", HandleHealth)

	return mux
}
