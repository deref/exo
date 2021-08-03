package container

import (
	"github.com/deref/exo/providers/docker"
	"github.com/deref/exo/providers/docker/compose"
)

type Container struct {
	docker.Component
	Spec
	State

	SyslogPort int
}

type Spec compose.Service

type State struct {
	ImageID     string `json:"imageId"`
	ContainerID string `json:"containerId"`
	Running     bool   `json:"running"`
}
