package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deref/exo/exod/api"
	"github.com/deref/exo/exod/plugin/process"
	"github.com/deref/exo/exod/statefile"
	"github.com/deref/pier"
)

var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
	mux.Handle("/_exo/", api.NewProviderMux("/_exo/", provider))
}

var provider = &Provider{}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, stateKey, (State)(statefile.New("state.json")))
	pier.Main(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mux.ServeHTTP(w, req.WithContext(ctx))
	}))
}

type contextKey int

const stateKey contextKey = 1

type Provider struct{}

func resolveProvider(typ string) (api.Provider, error) {
	switch typ {
	case "process":
		return &process.Provider{}, nil
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
