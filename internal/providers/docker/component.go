package docker

import (
	"strings"

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

func LoadSpec(s string, v compose.Interpolator, env map[string]string) error {
	if err := yamlutil.UnmarshalString(s, v); err != nil {
		return err
	}
	return v.Interpolate(compose.MapEnvironment(env))
}

func (c ComponentBase) LoadSpec(s string, v compose.Interpolator) error {
	return LoadSpec(s, v, c.WorkspaceEnvironment)
}

func (c ComponentBase) DockerObjectName(name, suffix string) string {
	var out strings.Builder
	if name == "" {
		out.WriteString(c.StackName)
		out.WriteByte('_')
		out.WriteString(c.ComponentName)
	} else {
		out.WriteString(name)
	}
	if suffix != "" {
		out.WriteByte('_')
		out.WriteString(suffix)
	}
	return out.String()
}
