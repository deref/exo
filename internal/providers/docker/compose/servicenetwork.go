package compose

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type ServiceNetworksTemplate []ServiceNetworkTemplate

type ServiceNetworkTemplate struct {
	IsShortForm bool
	Network     string
	ServiceNetworkTemplateLongForm
}

type ServiceNetworkTemplateLongForm struct {
	Aliases      []string `yaml:"aliases,omitempty"`
	IPv4Address  string   `yaml:"ipv4_address,omitempty"`
	IPv6Address  string   `yaml:"ipv6_address,omitempty"`
	LinkLocalIPs []string `yaml:"link_local_ips,omitempty"`
	Priority     string   `yaml:"priority,omitempty"`
}

type ServiceNetwork struct {
	IsShortForm bool
	Network     string
	ServiceNetworkLongForm
}

type ServiceNetworkLongForm struct {
	Aliases      []string `yaml:"aliases,omitempty"`
	IPV4Address  string   `yaml:"ipv4_address,omitempty"`
	IPV6Address  string   `yaml:"ipv6_address,omitempty"`
	LinkLocalIPs []string `yaml:"link_local_ips,omitempty"`
	Priority     int64    `yaml:"priority,omitempty"`
}

type ServiceNetworks []ServiceNetwork

func (sn ServiceNetworksTemplate) MarshalYAML() (interface{}, error) {
	slice := make(yaml.MapSlice, len(sn))
	for i, n := range sn {
		slice[i] = yaml.MapItem{
			Key:   n.Network,
			Value: n.ServiceNetworkTemplateLongForm,
		}
	}
	return slice, nil
}

func (sns *ServiceNetworksTemplate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asStrings []string
	if err := unmarshal(&asStrings); err == nil {
		nets := make([]ServiceNetworkTemplate, len(asStrings))
		for i, network := range asStrings {
			nets[i] = ServiceNetworkTemplate{
				IsShortForm: true,
				Network:     network,
			}
		}
		*sns = nets
		return nil
	}

	var mapSlice yaml.MapSlice
	if err := unmarshal(&mapSlice); err != nil {
		return err
	}

	var asMap map[string]ServiceNetworkTemplate
	if err := unmarshal(&asMap); err != nil {
		return err
	}

	nets := make([]ServiceNetworkTemplate, len(mapSlice))
	for i, item := range mapSlice {
		key, ok := item.Key.(string)
		if !ok {
			return fmt.Errorf("expected string key at index %d, got: %T", i, item.Key)
		}
		sn := asMap[key]
		sn.Network = key
		nets[i] = sn
	}
	*sns = nets

	return nil
}

func (sn ServiceNetworks) MarshalYAML() (interface{}, error) {
	slice := make(yaml.MapSlice, len(sn))
	for i, n := range sn {
		slice[i] = yaml.MapItem{
			Key:   n.Network,
			Value: n.ServiceNetworkLongForm,
		}
	}
	return slice, nil
}

func (sns *ServiceNetworks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asStrings []string
	if err := unmarshal(&asStrings); err == nil {
		nets := make([]ServiceNetwork, len(asStrings))
		for i, network := range asStrings {
			nets[i] = ServiceNetwork{
				Network: network,
			}
		}
		*sns = nets
		return nil
	}

	var mapSlice yaml.MapSlice
	if err := unmarshal(&mapSlice); err != nil {
		return err
	}

	var asMap map[string]ServiceNetwork
	if err := unmarshal(&asMap); err != nil {
		return err
	}

	nets := make([]ServiceNetwork, len(mapSlice))
	for i, item := range mapSlice {
		key, ok := item.Key.(string)
		if !ok {
			return fmt.Errorf("expected string key at index %d, got: %T", i, item.Key)
		}
		sn := asMap[key]
		sn.Network = key
		nets[i] = sn
	}
	*sns = nets

	return nil
}
