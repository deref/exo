// TODO: Generate this package via JOSH introspection.

package api

import (
	"context"
	"net/http"

	"github.com/deref/exo/exod/josh"
)

type Provider interface {
	Create(context.Context, *CreateInput) (*CreateOutput, error)
}

type CreateInput struct {
	Name string                 `json:"name"`
	Type string                 `json:"type"`
	Spec map[string]interface{} `json:"spec"`
}

type CreateOutput struct {
	IRI string `json:"iri"`
}

func NewProviderMux(prefix string, provider Provider) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(prefix+"create", josh.NewMethodHandler(provider.Create))
	return mux
}
