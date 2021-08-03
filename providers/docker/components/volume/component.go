package volume

import (
	"github.com/deref/exo/providers/docker"
	"github.com/deref/exo/providers/docker/compose"
)

type Volume struct {
	docker.Component
	Spec
	State
}

type Spec compose.Volume

type State struct {
	VolumeName string `json:"volumeId"`
}
