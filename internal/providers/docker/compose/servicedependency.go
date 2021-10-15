package compose

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ServiceDependencies struct {
	Style Style
	Items []ServiceDependency
}

type ServiceDependency struct {
	IsShortSyntax bool
	Service       string
	ServiceDependencyLongForm
}

type ServiceDependencyLongForm struct {
	Condition string `yaml:"condition,omitempty"`
}

func (deps *ServiceDependencies) UnmarshalYAML(node *yaml.Node) error {
	switch node.Tag {
	case "!!map":
		deps.Style = MapStyle
		n := len(node.Content) / 2
		deps.Items = make([]ServiceDependency, n)
		for i := 0; i < n; i++ {
			nameNode := node.Content[i*2+0]
			longFormNode := node.Content[i*2+1]
			var item ServiceDependency
			if err := nameNode.Decode(&item.Service); err != nil {
				return err
			}
			if err := longFormNode.Decode(&item.ServiceDependencyLongForm); err != nil {
				return err
			}
			deps.Items[i] = item
		}
		return nil
	case "!!seq":
		deps.Style = SeqStyle
		return node.Decode(&deps.Items)
	default:
		return fmt.Errorf("expected !!seq or !!map, got %s", node.ShortTag())
	}
}

func (deps ServiceDependencies) MarshalYAML() (interface{}, error) {
	if deps.Style == SeqStyle {
		return deps.Items, nil
	}
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, len(deps.Items)*2),
	}
	for i, item := range deps.Items {
		var keyNode, valueNode yaml.Node
		if err := keyNode.Encode(item.Service); err != nil {
			panic(err)
		}
		if err := valueNode.Encode(item); err != nil {
			return nil, err
		}
		node.Content[i*2+0] = &keyNode
		node.Content[i*2+1] = &valueNode
	}
	return node, nil
}

func (dep *ServiceDependency) UnmarshalYAML(node *yaml.Node) error {
	if node.Decode(&dep.Service) == nil {
		dep.IsShortSyntax = true
		dep.Condition = "service_started"
		return nil
	}
	return node.Decode(&dep.ServiceDependencyLongForm)
}

func (dep ServiceDependency) MarshalYAML() (interface{}, error) {
	if dep.IsShortSyntax && dep.Condition == "service_started" {
		return dep.Service, nil
	}
	return dep.ServiceDependencyLongForm, nil
}
