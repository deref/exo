package container

import (
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/providers/docker/compose"
)

type Container struct {
	docker.ComponentBase
	Spec
	State

	SyslogPort uint
}

type Spec compose.Service

type State struct {
	ImageID     string `json:"imageId"`
	ContainerID string `json:"containerId"`
	Running     bool   `json:"running"`
}
