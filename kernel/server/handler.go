package server

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/deref/exo/kernel/api"
	state "github.com/deref/exo/kernel/state/api"
	"github.com/deref/exo/kernel/state/statefile"
)

type Config struct {
	VarDir     string
	MuxPattern string
}

func NewContext(ctx context.Context, cfg *Config) context.Context {
	statePath := filepath.Join(cfg.VarDir, "state.json")
	return state.ContextWithStore(ctx, statefile.New(statePath))
}

func NewHandler(ctx context.Context, cfg *Config) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(cfg.MuxPattern, api.NewWorkspaceMux(cfg.MuxPattern, &Workspace{
		ID:     "default", // XXX
		VarDir: cfg.VarDir,
	}))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mux.ServeHTTP(w, req.WithContext(ctx))
	})
}
