package volume

import (
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/providers/docker/compose"
)

type Volume struct {
	docker.ComponentBase
	State
}

type Spec = compose.Volume

type State struct {
	VolumeName string `json:"volumeId"`
}
