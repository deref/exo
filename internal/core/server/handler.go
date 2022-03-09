package server

import (
	"net/http"

	"github.com/deref/exo/internal/about"
	"github.com/deref/exo/internal/api"
	state "github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/esv"
	"github.com/deref/exo/internal/task"
	"github.com/deref/exo/internal/token"
	"github.com/deref/exo/internal/util/logging"
	docker "github.com/docker/docker/client"
)

type Config struct {
	VarDir      string
	Store       state.Store
	Install     *about.Install
	SyslogPort  uint
	Docker      *docker.Client
	Logger      logging.Logger
	TaskTracker *task.TaskTracker
	TokenClient token.TokenClient
	EsvClient   esv.EsvClient
	ExoVersion  string
	Service     api.Service
}

type versionMiddleware struct {
	ExoVersion string
}

func (m *versionMiddleware) ServeHTTPMiddleware(w http.ResponseWriter, req *http.Request, next http.Handler) {
	w.Header().Add("Exo-Version", m.ExoVersion)
	next.ServeHTTP(w, req)
}

func BuildRootMux(prefix string, cfg *Config) *http.ServeMux {
	auth := &authMiddleware{
		TokenClient: cfg.TokenClient,
	}
	version := &versionMiddleware{
		ExoVersion: cfg.ExoVersion,
	}

	mux := http.NewServeMux()

	mux.Handle(prefix+"health", applyMiddleware(HandleHealth, version))

	mux.Handle(prefix+"graphql", applyMiddleware(&GraphqlHandler{
		Service: cfg.Service,
	}, version, auth))

	return mux
}
