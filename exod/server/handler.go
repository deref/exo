package server

import (
	"net/http"

	"github.com/deref/exo/exod/api"
	state "github.com/deref/exo/exod/state/api"
	josh "github.com/deref/exo/josh/server"
)

type Config struct {
	VarDir string
	Store  state.Store
}

func BuildRootMux(prefix string, cfg *Config) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)

	endKernel := b.Begin("kernel")
	api.BuildKernelMux(b, func(req *http.Request) api.Kernel {
		return &Kernel{
			VarDir: cfg.VarDir,
			Store:  cfg.Store,
		}
	})
	endKernel()

	endWorkspace := b.Begin("workspace")
	api.BuildWorkspaceMux(b, func(req *http.Request) api.Workspace {
		return &Workspace{
			ID:     req.URL.Query().Get("id"),
			VarDir: cfg.VarDir,
			Store:  cfg.Store,
		}
	})
	endWorkspace()

	return b.Build()
}
