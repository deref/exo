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
	Service       String
	ServiceDependencyLongForm
}

type ServiceDependencyLongForm struct {
	Condition String `yaml:"condition,omitempty"`
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

func (deps *ServiceDependencies) Interpolate(env Environment) error {
	return interpolateSlice(deps.Items, env)
}

func (deps ServiceDependencies) MarshalYAML() (any, error) {
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
	var err error
	if node.Tag == "!!str" {
		dep.IsShortSyntax = true
		err = node.Decode(&dep.Service)
	} else {
		err = node.Decode(&dep.ServiceDependencyLongForm)
	}
	_ = dep.Interpolate(ErrEnvironment)
	return err
}

func (dep *ServiceDependency) Interpolate(env Environment) error {
	if dep.IsShortSyntax {
		return dep.Service.Interpolate(env)
	}
	return dep.ServiceDependencyLongForm.Interpolate(env)
}

func (dep ServiceDependency) MarshalYAML() (any, error) {
	if dep.IsShortSyntax && dep.Condition.Value == "" {
		return dep.Service, nil
	}
	return dep.ServiceDependencyLongForm, nil
}

func (dep *ServiceDependencyLongForm) Interpolate(env Environment) error {
	return interpolateStruct(dep, env)
}
