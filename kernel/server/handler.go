package server

import (
	"context"
	"net/http"

	"github.com/deref/exo/kernel/api"
	"github.com/deref/exo/kernel/state"
	"github.com/deref/exo/kernel/state/statefile"
)

func NewContext(ctx context.Context) context.Context {
	return state.ContextWithStore(ctx, statefile.New("./var/state.json"))
}

var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
	mux.Handle("/", api.NewProjectMux("/", &Project{
		ID: "default", // XXX
	}))
}

func NewHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mux.ServeHTTP(w, req.WithContext(ctx))
	})
}
