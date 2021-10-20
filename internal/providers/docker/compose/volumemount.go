package compose

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type VolumeMount struct {
	ShortForm String
	VolumeMountLongForm
}

type VolumeMountLongForm struct {
	Type        String         `yaml:"type,omitempty"`
	Source      String         `yaml:"source,omitempty"`
	Target      String         `yaml:"target,omitempty"`
	ReadOnly    Bool           `yaml:"read_only,omitempty"`
	Bind        *BindOptions   `yaml:"bind,omitempty"`
	Volume      *VolumeOptions `yaml:"volume,omitempty"`
	Tmpfs       *TmpfsOptions  `yaml:"tmpfs,omitempty"`
	Consistency *Ignored       `yaml:"consistency,omitempty"`
}

type VolumeOptions struct {
	Nocopy Bool `yaml:"nocopy,omitempty"`
}

func (opt *VolumeOptions) Interpolate(env Environment) error {
	return interpolateStruct(opt, env)
}

type BindOptions struct {
	Propagation    String `yaml:"propagation,omitempty"`
	CreateHostPath Bool   `yaml:"create_host_path,omitempty"`
}

func (opt *BindOptions) Interpolate(env Environment) error {
	return interpolateStruct(opt, env)
}

type TmpfsOptions struct {
	Size Bytes `yaml:"size,omitempty"`
}

func (opt *TmpfsOptions) Interpolate(env Environment) error {
	return interpolateStruct(opt, env)
}

func (vm VolumeMount) MarshalYAML() (interface{}, error) {
	if vm.ShortForm.Expression != "" {
		return vm.ShortForm.Expression, nil
	}
	/* TODO: Inverse of interpolation goes elsewhere.
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
	*/
	return vm.VolumeMountLongForm, nil
}

func (vm *VolumeMount) UnmarshalYAML(node *yaml.Node) error {
	var err error
	if node.Tag == "!!str" {
		err = node.Decode(&vm.ShortForm)
	} else {
		err = node.Decode(&vm.VolumeMountLongForm)
	}
	_ = vm.Interpolate(ErrEnvironment)
	return err
}

func (vm *VolumeMount) Interpolate(env Environment) error {
	if vm.ShortForm.Expression == "" {
		return vm.VolumeMountLongForm.Interpolate(env)
	} else {
		if err := vm.ShortForm.Interpolate(env); err != nil {
			return nil
		}
		short := vm.ShortForm.Value
		parts := strings.Split(short, ":")
		switch len(parts) {
		case 1:
			vm.Type = MakeString("volume")
			vm.Target = vm.ShortForm
		case 2:
			vm.setSource(parts[0])
			vm.Target = MakeString(parts[1])
		case 3:
			vm.setSource(parts[0])
			vm.Target = MakeString(parts[1])
			accessMode := parts[2]
			switch accessMode {
			case "ro":
				vm.ReadOnly = MakeBool(true)
			case "rw":
				// Do nothing - va.ReadOnly is already false.
			case "cached", "delegated":
				// Legacy read/write modes that no longer have any effect.
			default:
				return fmt.Errorf(`invalid access mode; expected "ro" or "rw" but got %q`, accessMode)
			}
		default:
			return fmt.Errorf(`invalid volume specification; expected "VOLUME:CONTAINER_PATH" or "VOLUME:CONTAINER_PATH:ACCESS_MODE" but got %q`, short)
		}
		return nil
	}
}

func (vm *VolumeMountLongForm) Interpolate(env Environment) error {
	return interpolateStruct(vm, env)
}

var localPathRe = regexp.MustCompile("^[./~]")

func (vm *VolumeMount) setSource(src string) {
	vm.Source = MakeString(src)
	if localPathRe.MatchString(src) {
		vm.Type = MakeString("bind")
		vm.Bind = &BindOptions{}
		// CreateHostPath is always implied by the short syntax.
		vm.Bind.CreateHostPath = MakeBool(true)
	} else {
		vm.Type = MakeString("volume")
	}
}
