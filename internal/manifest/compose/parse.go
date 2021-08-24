package compose

import (
	"io"

	"github.com/goccy/go-yaml"
)

// TODO: Consider whether this representation is useful or if we can just use `compose.Compose`,
type Project struct {
	Services map[string]yaml.MapSlice `yaml:"services"`
	Networks map[string]yaml.MapSlice `yaml:"networks"`
	Volumes  map[string]yaml.MapSlice `yaml:"volumes"`
}

func Parse(r io.Reader) (*Project, error) {
	// TODO: Preserve comments. This is supported by a later version of go-yaml.
	dec := yaml.NewDecoder(r, yaml.DisallowDuplicateKey())
	var project Project
	if err := dec.Decode(&project); err != nil {
		return nil, err
	}
	return &project, nil
}
