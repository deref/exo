package server

import (
	"net/http"

	"github.com/deref/exo/core/api"
	state "github.com/deref/exo/core/state/api"
	josh "github.com/deref/exo/josh/server"
	"github.com/deref/exo/telemetry"
	"github.com/deref/exo/util/logging"
	docker "github.com/docker/docker/client"
)

type Config struct {
	VarDir     string
	Store      state.Store
	Telemetry  telemetry.Telemetry
	SyslogPort int
	Docker     *docker.Client
	Logger     logging.Logger
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
			Logger:     cfg.Logger,
			Store:      cfg.Store,
			SyslogPort: cfg.SyslogPort,
			Docker:     cfg.Docker,
		}
	})
	endWorkspace()

	mux := b.Build()

	mux.Handle(prefix+"health", HandleHealth)

	return mux
}
