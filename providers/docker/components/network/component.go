package network

import (
	"github.com/deref/exo/providers/docker"
	"github.com/deref/exo/providers/docker/compose"
)

type Network struct {
	docker.Component
	Spec
	State
}

type Spec compose.Network

type State struct {
	NetworkID string `json:"networkId"`
}
