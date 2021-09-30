package network

import (
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/providers/docker/compose"
)

type Network struct {
	docker.ComponentBase
	State
}

type Spec compose.Network

type State struct {
	NetworkID string `json:"networkId"`
}
