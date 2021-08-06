package server

import (
	"net/http"

	"github.com/deref/exo/internal/core/api"
	state "github.com/deref/exo/internal/core/state/api"
	josh "github.com/deref/exo/internal/josh/server"
	"github.com/deref/exo/internal/task"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/logging"
	docker "github.com/docker/docker/client"
)

type Config struct {
	VarDir      string
	Store       state.Store
	Telemetry   telemetry.Telemetry
	SyslogPort  uint
	Docker      *docker.Client
	Logger      logging.Logger
	TaskTracker *task.TaskTracker
}

func BuildRootMux(prefix string, cfg *Config) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)

	endKernel := b.Begin("kernel")
	api.BuildKernelMux(b, func(req *http.Request) api.Kernel {
		return &Kernel{
			VarDir:      cfg.VarDir,
			Store:       cfg.Store,
			Telemetry:   cfg.Telemetry,
			TaskTracker: cfg.TaskTracker,
		}
	})
	endKernel()

	endWorkspace := b.Begin("workspace")
	api.BuildWorkspaceMux(b, func(req *http.Request) api.Workspace {
		return &Workspace{
			ID:          req.URL.Query().Get("id"),
			VarDir:      cfg.VarDir,
			Logger:      cfg.Logger,
			Store:       cfg.Store,
			SyslogPort:  cfg.SyslogPort,
			Docker:      cfg.Docker,
			TaskTracker: cfg.TaskTracker,
		}
	})
	endWorkspace()

	endTaskStore := b.Begin("task-store")
	taskapi.BuildTaskStoreMux(b, func(req *http.Request) taskapi.TaskStore {
		return cfg.TaskTracker.Store
	})
	endTaskStore()

	mux := b.Build()

	mux.Handle(prefix+"health", HandleHealth)

	return mux
}
