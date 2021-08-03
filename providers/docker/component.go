package docker

import (
	"github.com/deref/exo/providers/core"
	dockerclient "github.com/docker/docker/client"
)

type Component struct {
	core.Component
	Docker *dockerclient.Client
}
