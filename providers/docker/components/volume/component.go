package volume

import (
	docker "github.com/docker/docker/client"
)

type Volume struct {
	ComponentID string
	Spec
	State

	Docker *docker.Client
}

// See note: [COMPOSE_YAML].
type Spec struct {
	Driver     string            `yaml:"driver"`
	DriverOpts map[string]string `yaml:"driver_opts"`
	// TODO: external
	Labels map[string]string `yaml:"labels"` // TODO: Support array syntax.
	Name   string            `yaml:"name"`
}

type State struct {
	VolumeName string `json:"volumeId"`
}
