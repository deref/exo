package compose

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

type VolumeMount struct {
	Type        string
	Source      string
	Target      string
	ReadOnly    bool
	Bind        *BindOptions
	Volume      *VolumeOptions
	Tmpfs       *TmpfsOptions
	Consistency *IgnoredField
}

// extendedVolumeMount is a private struct that is structurally identical to VolumeAttachment but
// is only used for YAML unmarshalling where we do not need to consider the short string-based syntax.
type extendedVolumeMount struct {
	Type        string         `yaml:"type"`
	Source      string         `yaml:"source"`
	Target      string         `yaml:"target"`
	ReadOnly    bool           `yaml:"read_only"`
	Bind        *BindOptions   `yaml:"bind"`
	Volume      *VolumeOptions `yaml:"volume"`
	Tmpfs       *TmpfsOptions  `yaml:"tmpfs"`
	Consistency *IgnoredField  `yaml:"consistency"`
}

func (vm *VolumeMount) UnmarshalYAML(b []byte) error {
	var asString string
	if err := yaml.Unmarshal(b, &asString); err == nil {
		return vm.fromShortSyntax(asString)
	}

	asExtended := extendedVolumeMount{}
	if err := yaml.Unmarshal(b, &asExtended); err != nil {
		return err
	}
	vm.Type = asExtended.Type
	vm.Source = asExtended.Source
	vm.Target = asExtended.Target
	vm.ReadOnly = asExtended.ReadOnly
	vm.Bind = asExtended.Bind
	vm.Volume = asExtended.Volume
	vm.Tmpfs = asExtended.Tmpfs
	vm.Consistency = asExtended.Consistency

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
	Nocopy bool `yaml:"nocopy"`
}

type BindOptions struct {
	Propagation    string `yaml:"propagation"`
	CreateHostPath bool   `yaml:"create_host_path"`
}

type TmpfsOptions struct {
	Size int64 `yaml:"size"`
}