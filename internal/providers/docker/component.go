package docker

import (
	"github.com/deref/exo/internal/providers/core"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/deref/exo/internal/util/yamlutil"
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

func (c *ComponentBase) UnmarshalSpec(spec string, v interface{}) error {
	return yamlutil.UnmarshalString(spec, &compose.Interpolated{
		Environment: compose.MapEnvironment(c.WorkspaceEnvironment),
		Value:       v,
	})
}
