package compose

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type VolumeMount struct {
	IsShortForm bool
	VolumeMountLongForm
}

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

func (vm VolumeMount) MarshalYAML() (interface{}, error) {
	if vm.IsShortForm {
		var sb strings.Builder
		if vm.Source != "" {
			sb.WriteString(vm.Source)
			sb.WriteRune(':')
		}
		sb.WriteString(vm.Target)
		if vm.ReadOnly {
			sb.WriteString(":ro")
		}
		return sb.String(), nil
	}
	return vm.VolumeMountLongForm, nil
}

func (vm *VolumeMount) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if node.Decode(&s) == nil {
		vm.IsShortForm = true
		return vm.fromShortSyntax(s)
	}
	return node.Decode(&vm.VolumeMountLongForm)
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
