package apigateway

import (
	"github.com/deref/exo/internal/providers/docker/components/container"
	"github.com/deref/exo/internal/token"
)

type Spec struct {
	Port int `json:"port"`
}

type APIGateway struct {
	TokenClient token.TokenClient
	container.Container
}
