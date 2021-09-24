package server

import (
	"net/http"
	"strings"

	"github.com/deref/exo/internal/core/api"
	state "github.com/deref/exo/internal/core/state/api"
	josh "github.com/deref/exo/internal/josh/server"
	"github.com/deref/exo/internal/task"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/token"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/httputil"
	"github.com/deref/exo/internal/util/logging"
	docker "github.com/docker/docker/client"
)

type Config struct {
	VarDir      string
	Store       state.Store
	SyslogPort  uint
	Docker      *docker.Client
	Logger      logging.Logger
	TaskTracker *task.TaskTracker
	TokenClient token.TokenClient
}

func BuildRootMux(prefix string, cfg *Config) *http.ServeMux {
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
			cookie, err := req.Cookie("token")
			if token == "" && err == nil {
				token = cookie.Value
			}

			authed, err := cfg.TokenClient.CheckToken(token)
			if err != nil {
				httputil.WriteError(w, req, errutil.NewHTTPError(http.StatusInternalServerError, "Could not validate token"))
				return
			}
			if !authed {
				httputil.WriteError(w, req, errutil.NewHTTPError(http.StatusUnauthorized, "Bad or no token"))
				return
			}
			next.ServeHTTP(w, req)
		})
	}

	b := josh.NewMuxBuilder(prefix)
	b.AddMiddleware(authMiddleware)

	endKernel := b.Begin("kernel")
	api.BuildKernelMux(b, func(req *http.Request) api.Kernel {
		return &Kernel{
			VarDir:      cfg.VarDir,
			Store:       cfg.Store,
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
