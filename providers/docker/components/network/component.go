package network

import (
	docker "github.com/docker/docker/client"
)

type Network struct {
	ComponentID string
	Spec
	State

	Docker *docker.Client
}

// See note: [COMPOSE_YAML].
type Spec struct {
	Driver     string            `yaml:"driver"`
	DriverOpts map[string]string `yaml:"driver_opts"`
	Attachable bool              `yaml:"attachable"`
	EnableIPv6 bool              `yaml:"enable_ipv6"`
	Internal   bool              `yaml:"internal"`
	Labels     map[string]string `yaml:"labels"` // TODO: Support array syntax.
	External   bool              `yaml:"external"`
	// TODO: name
}

type State struct {
	NetworkID string `json:"networkId"`
}
