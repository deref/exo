package compose

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Ulimits []Ulimit

type Ulimit struct {
	Name string

	ShortForm Int
	UlimitLongForm
}

type UlimitLongForm struct {
	Soft Int `yaml:"soft,omitempty"`
	Hard Int `yaml:"hard,omitempty"`
}

func (uls *Ulimits) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag != "!!map" {
		return fmt.Errorf("expected !!map, got %s", node.ShortTag())
	}
	n := len(node.Content) / 2
	items := make([]Ulimit, n)
	for i := 0; i < n; i++ {
		nameNode := node.Content[i*2+0]
		valueNode := node.Content[i*2+1]
		if err := valueNode.Decode(&items[i]); err != nil {
			return err
		}
		if err := nameNode.Decode(&items[i].Name); err != nil {
			return err
		}
	}
	*uls = items
	return nil
}

func (uls *Ulimits) Interpolate(env Environment) error {
	return interpolateSlice(*uls, env)
}

func (uls Ulimits) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, len(uls)*2),
	}
	for i, item := range uls {
		var keyNode, valueNode yaml.Node
		if err := keyNode.Encode(item.Name); err != nil {
			return nil, err
		}
		if err := valueNode.Encode(item); err != nil {
			return nil, err
		}
		node.Content[i*2+0] = &keyNode
		node.Content[i*2+1] = &valueNode
	}
	return node, nil
}

func (ul *Ulimit) UnmarshalYAML(node *yaml.Node) error {
	var err error
	if node.Kind == yaml.ScalarNode {
		err = node.Decode(&ul.ShortForm)
	} else {
		err = node.Decode(&ul.UlimitLongForm)
	}
	_ = ul.Interpolate(ErrEnvironment)
	return err
}

func (ul *Ulimit) Interpolate(env Environment) error {
	if ul.ShortForm.Tag != "" {
		if err := ul.ShortForm.Interpolate(env); err != nil {
			return err
		}
		ul.Soft = ul.ShortForm
		ul.Hard = ul.ShortForm
	}
	return ul.UlimitLongForm.Interpolate(env)
}

func (ul Ulimit) MarshalYAML() (interface{}, error) {
	if ul.ShortForm.Tag != "" {
		return ul.ShortForm, nil
	}
	return ul.UlimitLongForm, nil
}

func (ul *UlimitLongForm) Interpolate(env Environment) error {
	return interpolateStruct(ul, env)
}
