package docker

import (
	"github.com/deref/exo/internal/providers/core"
	dockerclient "github.com/docker/docker/client"
)

type ComponentBase struct {
	core.ComponentBase
	Docker *dockerclient.Client
}

func (c ComponentBase) GetExoLabels() map[string]string {
	return map[string]string{
		"io.deref.exo.workspace": c.WorkspaceID,
		"io.deref.exo.component": c.ComponentID,
	}
}
