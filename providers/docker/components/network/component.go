package network

import (
	"github.com/deref/exo/providers/docker/compose"
	docker "github.com/docker/docker/client"
)

type Network struct {
	ComponentID string
	Spec
	State

	Docker *docker.Client
}

type Spec compose.Network

type State struct {
	NetworkID string `json:"networkId"`
}
