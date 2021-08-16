package container

import (
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/providers/docker/compose"
)

type Container struct {
	docker.ComponentBase

	Spec  Spec
	State State

	SyslogPort uint
}

type Spec compose.Service

type State struct {
	ContainerID string     `json:"containerId"`
	Running     bool       `json:"running"`
	Image       ImageState `json:"image"`
}

type ImageState struct {
	ID         string   `json:"id"`
	Command    []string `json:"command"`
	WorkingDir string   `json:"workingDir"`
	Entrypoint []string `json:"entrypoint"`
}
