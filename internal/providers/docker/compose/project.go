// References:
// https://github.com/compose-spec/compose-spec/blob/master/spec.md
// https://docs.docker.com/compose/compose-file/compose-file-v3/
// https://github.com/docker/compose/blob/4a51af09d6cdb9407a6717334333900327bc9302/compose/config/compose_spec.json
//
// Fields enumerated as of July 17, 2021 with from the following spec file:
// <https://github.com/compose-spec/compose-spec/blob/5141aafafa6ea03fcf52eb2b44218408825ab480/spec.md>.

package compose

import (
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-yaml"
)

func Parse(r io.Reader) (*Project, error) {
	dec := yaml.NewDecoder(r,
		yaml.DisallowDuplicateKey(),
	)
	var comp Project
	if err := dec.Decode(&comp); err != nil {
		return nil, err
	}

	// Validate.
	for key := range comp.Raw {
		switch key {
		case "version", "services", "networks", "volumes", "configs", "secrets":
			// Ok.
		default:
			if !strings.HasPrefix(key, "x-") {
				return nil, fmt.Errorf("unsupported top-level key in compose file: %q", key)
			}
		}
	}

	return &comp, nil
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
}
