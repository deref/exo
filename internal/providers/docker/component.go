package docker

import (
	"github.com/deref/exo/internal/providers/core"
	dockerclient "github.com/docker/docker/client"
)

type ComponentBase struct {
	core.ComponentBase
	Docker *dockerclient.Client
}
