package docker

import (
	"fmt"

	"github.com/deref/exo/internal/providers/core"
	"github.com/deref/exo/internal/providers/docker/compose"
	dockerclient "github.com/docker/docker/client"
	"github.com/goccy/go-yaml"
)

type ComponentBase struct {
	core.ComponentBase
	Docker *dockerclient.Client
}

func LoadSpec(spec string, v interface{}, env map[string]string) error {
	if err := yaml.Unmarshal([]byte(spec), v); err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	if err := compose.Interpolate(v, compose.MapEnvironment(env)); err != nil {
		return fmt.Errorf("interpolating: %w", err)
	}
	return nil
}
