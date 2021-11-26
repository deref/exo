package apigateway

import (
	"github.com/deref/exo/internal/providers/docker/components/container"
	"github.com/deref/exo/internal/token"
)

type Spec struct {
	WebPort int `json:"web_port"`
	APIPort int `json:"api_port"`
}

type State struct {
	WebPort int `json:"web_port"`
	APIPort int `json:"api_port"`
	container.State
}

type APIGateway struct {
	TokenClient token.TokenClient
	container.Container
	State State
}
