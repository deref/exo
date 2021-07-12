package exo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deref/exo/api"
	"github.com/deref/exo/plugin/process"
	"github.com/deref/exo/statefile"
)

func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, stateKey, (State)(statefile.New("state.json")))
}

var mux *http.ServeMux

func NewHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mux.ServeHTTP(w, req.WithContext(ctx))
	})
}

func init() {
	mux = http.NewServeMux()
	mux.Handle("/_exo/", api.NewProviderMux("/_exo/", provider))
}

var provider = &Provider{}

type contextKey int

const stateKey contextKey = 1

type Provider struct{}

func resolveProvider(typ string) (api.Provider, error) {
	switch typ {
	case "process":
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		projectDir := wd                   // TODO: Get from component hierarchy.
		varDir := filepath.Join(wd, "var") // TODO: Get from exod config.
		return &process.Provider{
			ProjectDir: projectDir,
			VarDir:     filepath.Join(varDir, "proc"),
		}, nil
	}
	return nil, fmt.Errorf("no provider for type: %q", typ)
}

func (provider *Provider) Create(ctx context.Context, input *api.CreateInput) (*api.CreateOutput, error) {
	underlying, err := resolveProvider(input.Type)
	if err != nil {
		return nil, err
	}

	state := ctx.Value(stateKey).(State)

	if !IsValidName(input.Name) {
		return nil, fmt.Errorf("invalid name: %q", input.Name)
	}

	// TODO: Record pending create, to aid in recover on failure.

	output, err := underlying.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := state.Remember(ctx, input.Name, output.IRI); err != nil {
		return nil, fmt.Errorf("remembering created component: %w", err)
	}

	return output, nil
}

func IsValidName(name string) bool {
	return name != "" // XXX
}

type State interface {
	Remember(ctx context.Context, name string, iri string) error
}
