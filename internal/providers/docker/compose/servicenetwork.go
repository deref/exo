package compose

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ServiceNetworks struct {
	Style Style
	Items []ServiceNetwork
}

type ServiceNetwork struct {
	Key string

	ShortForm String
	ServiceNetworkLongForm
}

type ServiceNetworkLongForm struct {
	Aliases      Strings `yaml:"aliases,omitempty"`
	IPV4Address  String  `yaml:"ipv4_address,omitempty"`
	IPV6Address  String  `yaml:"ipv6_address,omitempty"`
	LinkLocalIPs Strings `yaml:"link_local_ips,omitempty"`
	Priority     Int     `yaml:"priority,omitempty"`
}

func (sn *ServiceNetworks) UnmarshalYAML(node *yaml.Node) error {
	switch node.Tag {
	case "!!map":
		sn.Style = MapStyle
		n := len(node.Content) / 2
		sn.Items = make([]ServiceNetwork, n)
		for i := 0; i < n; i++ {
			keyNode := node.Content[i*2+0]
			longFormNode := node.Content[i*2+1]
			var item ServiceNetwork
			if err := keyNode.Decode(&item.Key); err != nil {
				return err
			}
			if err := longFormNode.Decode(&item.ServiceNetworkLongForm); err != nil {
				return err
			}
			sn.Items[i] = item
		}
		return nil
	case "!!seq":
		sn.Style = SeqStyle
		return node.Decode(&sn.Items)
	default:
		return fmt.Errorf("expected !!seq or !!map, got %s", node.ShortTag())
	}
}

func (sn *ServiceNetworks) Interpolate(env Environment) error {
	return interpolateSlice(sn.Items, env)
}

func (sn ServiceNetworks) MarshalYAML() (any, error) {
	if sn.Style == UnknownStyle {
		sn.Style = SeqStyle
		for _, item := range sn.Items {
			if item.ShortForm.Expression == "" {
				sn.Style = MapStyle
			}
		}
	}
	if sn.Style == SeqStyle {
		return sn.Items, nil
	}
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, len(sn.Items)*2),
	}
	for i, item := range sn.Items {
		var keyNode, valueNode yaml.Node
		if err := keyNode.Encode(item.Key); err != nil {
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

func (sn *ServiceNetwork) UnmarshalYAML(node *yaml.Node) error {
	var err error
	if node.Tag == "!!str" {
		err = node.Decode(&sn.ShortForm)
	} else {
		err = node.Decode(&sn.ServiceNetworkLongForm)
	}
	_ = sn.Interpolate(ErrEnvironment)
	return err
}

func (sn *ServiceNetwork) Interpolate(env Environment) error {
	if sn.ShortForm.Tag != "" {
		if err := sn.ShortForm.Interpolate(env); err != nil {
			return err
		}
		sn.Key = sn.ShortForm.Value
	}
	return sn.ServiceNetworkLongForm.Interpolate(env)
}

func (sn ServiceNetwork) MarshalYAML() (any, error) {
	if sn.ShortForm.Expression != "" {
		return sn.ShortForm.Expression, nil
	}
	return sn.ServiceNetworkLongForm, nil
}

func (sn *ServiceNetworkLongForm) Interpolate(env Environment) error {
	return interpolateStruct(sn, env)
}
