package image

import (
	"github.com/deref/exo/internal/providers/docker/compose"
)

// SEE NOTE [IMAGE_SUBCOMPONENT].

type Spec struct {
	Platform string        `yaml:"platform"`
	Build    compose.Build `yaml:"build"`
}

func (s *Spec) Interpolate(env compose.Environment) error {
	return s.Build.Interpolate(env)
}
