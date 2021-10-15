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
	Key         string
	IsShortForm bool
	ServiceNetworkLongForm
}

type ServiceNetworkLongForm struct {
	Aliases      []string `yaml:"aliases,omitempty"`
	IPV4Address  string   `yaml:"ipv4_address,omitempty"`
	IPV6Address  string   `yaml:"ipv6_address,omitempty"`
	LinkLocalIPs []string `yaml:"link_local_ips,omitempty"`
	Priority     int64    `yaml:"priority,omitempty"`
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

func (sn ServiceNetworks) MarshalYAML() (interface{}, error) {
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
	if node.Decode(&sn.Key) == nil {
		sn.IsShortForm = true
		return nil
	}
	return node.Decode(&sn.ServiceNetworkLongForm)
}

func (sn ServiceNetwork) MarshalYAML() (interface{}, error) {
	if sn.IsShortForm {
		return sn.Key, nil
	}
	return sn.ServiceNetworkLongForm, nil
}
