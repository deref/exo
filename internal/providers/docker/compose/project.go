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
	"reflect"

	"gopkg.in/yaml.v3"
)

func Parse(r io.Reader) (*Project, error) {
	dec := yaml.NewDecoder(r)
	var comp Project
	if err := dec.Decode(&comp); err != nil {
		return nil, err
	}
	// TODO: Initial validation pass.  This pass has to be separate, since it
	// needs to be possible to re-run it after string interpolation.
	return &comp, nil
}

type Project struct {
	Version  String          `yaml:"version,omitempty"`
	Services ProjectServices `yaml:"services,omitempty"`
	Networks ProjectNetworks `yaml:"networks,omitempty"`
	Volumes  ProjectVolumes  `yaml:"volumes,omitempty"`
	Configs  ProjectConfigs  `yaml:"configs,omitempty"`
	Secrets  ProjectSecrets  `yaml:"secrets,omitempty"`
}

type ProjectServices []Service
type ProjectNetworks []Network
type ProjectVolumes []Volume
type ProjectConfigs []Config
type ProjectSecrets []Secret

func (project *Project) Interpolate(env Environment) error {
	return interpolateStruct(project, env)
}

func (section *ProjectServices) UnmarshalYAML(node *yaml.Node) error {
	return unmarshalSection(section, node)
}
func (section *ProjectNetworks) UnmarshalYAML(node *yaml.Node) error {
	return unmarshalSection(section, node)
}
func (section *ProjectVolumes) UnmarshalYAML(node *yaml.Node) error {
	return unmarshalSection(section, node)
}
func (section *ProjectConfigs) UnmarshalYAML(node *yaml.Node) error {
	return unmarshalSection(section, node)
}
func (section *ProjectSecrets) UnmarshalYAML(node *yaml.Node) error {
	return unmarshalSection(section, node)
}

func (section ProjectServices) MarshalYAML() (any, error) {
	return marshalSection(section)
}
func (section ProjectNetworks) MarshalYAML() (any, error) {
	return marshalSection(section)
}
func (section ProjectVolumes) MarshalYAML() (any, error) {
	return marshalSection(section)
}
func (section ProjectConfigs) MarshalYAML() (any, error) {
	return marshalSection(section)
}
func (section ProjectSecrets) MarshalYAML() (any, error) {
	return marshalSection(section)
}

func (section *ProjectServices) Interpolate(env Environment) error {
	return interpolateSlice(*section, env)
}
func (section *ProjectNetworks) Interpolate(env Environment) error {
	return interpolateSlice(*section, env)
}
func (section *ProjectVolumes) Interpolate(env Environment) error {
	return interpolateSlice(*section, env)
}
func (section *ProjectConfigs) Interpolate(env Environment) error {
	return interpolateSlice(*section, env)
}
func (section *ProjectSecrets) Interpolate(env Environment) error {
	return interpolateSlice(*section, env)
}

func unmarshalSection(v any, node *yaml.Node) error {
	if node.Tag != "!!map" {
		return fmt.Errorf("expected !!map, got %s", node.Tag)
	}
	rv := reflect.ValueOf(v)
	n := len(node.Content) / 2
	rv.Elem().Set(reflect.MakeSlice(rv.Type().Elem(), n, n))
	for i := 0; i < n; i++ {
		keyNode := node.Content[i*2+0]
		configNode := node.Content[i*2+1]
		elem := rv.Elem().Index(i)
		if err := configNode.Decode(elem.Addr().Interface()); err != nil {
			return err
		}
		if err := keyNode.Decode(elem.FieldByName("Key").Addr().Interface()); err != nil {
			return err
		}
	}
	return nil
}

func marshalSection(v any) (any, error) {
	rv := reflect.ValueOf(v)
	n := rv.Len()
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, n*2),
	}
	for i := 0; i < n; i++ {
		elem := rv.Index(i)
		var keyNode, valueNode yaml.Node
		if err := keyNode.Encode(elem.FieldByName("Key").Interface()); err != nil {
			panic(err)
		}
		if err := valueNode.Encode(elem.Interface()); err != nil {
			return nil, err
		}
		node.Content[i*2+0] = &keyNode
		node.Content[i*2+1] = &valueNode
	}
	return node, nil
}
