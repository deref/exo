// References:
// https://github.com/compose-spec/compose-spec/blob/master/spec.md
// https://docs.docker.com/compose/compose-file/compose-file-v3/
// https://github.com/docker/compose/blob/4a51af09d6cdb9407a6717334333900327bc9302/compose/config/compose_spec.json
//
// Fields enumerated as of July 17, 2021 with from the following spec file:
// <https://github.com/compose-spec/compose-spec/blob/5141aafafa6ea03fcf52eb2b44218408825ab480/spec.md>.

package compose

import (
	"github.com/goccy/go-yaml"
)

type ProjectTemplate struct {
	Services map[string]ServiceTemplate `yaml:"services,omitempty"`
	Volumes  map[string]VolumeTemplate  `yaml:"volumes,omitempty"`
	Networks map[string]NetworkTemplate `yaml:"networks,omitempty"`
	MapSlice yaml.MapSlice              `yaml:",inline"`
}

type Project struct {
	Version  string             `yaml:"version,omitempty"`
	Services map[string]Service `yaml:"services,omitempty"`
	Networks map[string]Network `yaml:"networks,omitempty"`
	Volumes  map[string]Volume  `yaml:"volumes,omitempty"`
	Configs  map[string]Config  `yaml:"configs,omitempty"`
	Secrets  map[string]Secret  `yaml:"secrets,omitempty"`

	Raw map[string]interface{} `yaml:",inline"`
	// TODO: extensions with "x-" prefix.
	// TODO: Validate set of top-level keys.
}
