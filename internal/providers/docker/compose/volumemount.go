package compose

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

type VolumeMountTemplate struct {
	IsShortForm bool
	VolumeMountTemplateLongForm
}

type VolumeMountTemplateLongForm struct {
	Type     string        `yaml:"type,omitempty"`
	Source   string        `yaml:"source,omitempty"`
	MapSlice yaml.MapSlice `yaml:",inline"`
}

type VolumeMount struct {
	IsShortForm bool
	VolumeMountLongForm
}

// SEE NOTE [COMPOSE_AST].
type VolumeMountLongForm struct {
	Type        string         `yaml:"type,omitempty"`
	Source      string         `yaml:"source,omitempty"`
	Target      string         `yaml:"target,omitempty"`
	ReadOnly    bool           `yaml:"read_only,omitempty"`
	Bind        *BindOptions   `yaml:"bind,omitempty"`
	Volume      *VolumeOptions `yaml:"volume,omitempty"`
	Tmpfs       *TmpfsOptions  `yaml:"tmpfs,omitempty"`
	Consistency *Ignored       `yaml:"consistency,omitempty"`
}

func (vm VolumeMountTemplate) MarshalYAML() (interface{}, error) {
	return vm.VolumeMountTemplateLongForm, nil
}

func (vm *VolumeMountTemplate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asString string
	if err := unmarshal(&asString); err == nil {
		return vm.fromShortSyntax(asString)
	}
	return unmarshal(&vm.VolumeMountTemplateLongForm)
}

func (vm *VolumeMount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asString string
	if err := unmarshal(&asString); err == nil {
		return vm.fromShortSyntax(asString)
	}
	return unmarshal(&vm.VolumeMountLongForm)
}

func (vm *VolumeMountTemplate) fromShortSyntax(in string) error {
	parts := strings.Split(in, ":")
	switch len(parts) {
	case 1:
		vm.Type = "volume"
	case 2, 3:
		vm.Source = parts[0]
		if localPathRe.MatchString(vm.Source) {
			vm.Type = "bind"
		} else {
			vm.Type = "volume"
		}
	default:
		return fmt.Errorf(`invalid volume specification; expected "VOLUME:CONTAINER_PATH" or "VOLUME:CONTAINER_PATH:ACCESS_MODE" but got %q`, in)
	}
	return nil
}

func (vm *VolumeMount) fromShortSyntax(in string) error {
	parts := strings.Split(in, ":")
	switch len(parts) {
	case 1:
		vm.Type = "volume"
		vm.Target = in
	case 2:
		vm.setSource(parts[0])
		vm.Target = parts[1]
	case 3:
		vm.setSource(parts[0])
		vm.Target = parts[1]
		accessMode := parts[2]
		switch accessMode {
		case "ro":
			vm.ReadOnly = true
		case "rw":
			// Do nothing - va.ReadOnly is already false.
		case "cached", "delegated":
			// Legacy read/write modes that no longer have any effect.
		default:
			return fmt.Errorf(`invalid access mode; expected "ro" or "rw" but got %q`, accessMode)
		}
	default:
		return fmt.Errorf(`invalid volume specification; expected "VOLUME:CONTAINER_PATH" or "VOLUME:CONTAINER_PATH:ACCESS_MODE" but got %q`, in)
	}

	return nil
}

var localPathRe = regexp.MustCompile("^[./~]")

func (vm *VolumeMount) setSource(src string) {
	vm.Source = src
	if localPathRe.MatchString(src) {
		vm.Type = "bind"
		vm.Bind = &BindOptions{
			// CreateHostPath is always implied by the short syntax.
			CreateHostPath: true,
		}
	} else {
		vm.Type = "volume"
	}
}

type VolumeOptions struct {
	Nocopy bool `yaml:"nocopy,omitempty"`
}

type BindOptions struct {
	Propagation    string `yaml:"propagation,omitempty"`
	CreateHostPath bool   `yaml:"create_host_path,omitempty"`
}

type TmpfsOptions struct {
	Size int64 `yaml:"size,omitempty"`
}
