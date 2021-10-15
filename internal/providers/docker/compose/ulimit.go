package compose

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Ulimits []Ulimit

type Ulimit struct {
	Name        string
	IsShortForm bool
	UlimitLongForm
}

type UlimitLongForm struct {
	Soft int64 `yaml:"soft,omitempty"`
	Hard int64 `yaml:"hard,omitempty"`
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
	var n int64
	if node.Decode(&n) == nil {
		ul.IsShortForm = true
		ul.Soft = n
		ul.Hard = n
		return nil
	}
	return node.Decode(&ul.UlimitLongForm)
}

func (ul Ulimit) MarshalYAML() (interface{}, error) {
	if ul.IsShortForm && ul.Hard == ul.Soft {
		return ul.Soft, nil
	}
	return ul.UlimitLongForm, nil
}
