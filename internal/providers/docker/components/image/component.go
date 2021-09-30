package image

import (
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/deref/exo/internal/util/yamlutil"
)

// SEE NOTE [IMAGE_SUBCOMPONENT].

type Spec struct {
	Platform string        `yaml:"platform"`
	Build    compose.Build `yaml:"build"`
}

func UnmarshalSpec(s string) (*Spec, error) {
	var spec Spec
	err := yamlutil.UnmarshalString(s, &spec)
	return &spec, err
}
