package volume

import (
	"github.com/deref/exo/providers/docker/compose"
	docker "github.com/docker/docker/client"
)

type Volume struct {
	ComponentID string
	Spec
	State

	Docker *docker.Client
}

type Spec compose.Volume

type State struct {
	VolumeName string `json:"volumeId"`
}
