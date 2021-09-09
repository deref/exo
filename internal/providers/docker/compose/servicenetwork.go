package compose

import (
	"github.com/goccy/go-yaml"
)

type ServiceNetwork struct {
	Network      string
	Aliases      []string
	IPV4Address  string
	IPV6Address  string
	LinkLocalIPs []string
	Priority     int64
}

type ServiceNetworks []ServiceNetwork

func (sn *ServiceNetworks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asStrings []string
	if err := unmarshal(&asStrings); err == nil {
		nets := make([]ServiceNetwork, len(asStrings))
		for i, network := range asStrings {
			nets[i] = ServiceNetwork{
				Network: network,
			}
		}
		*sn = nets
		return nil
	}

	var asMap yaml.MapSlice
	if err := unmarshal(&asMap); err != nil {
		return err
	}

	nets := make([]ServiceNetwork, len(asMap))
	for i, item := range asMap {
		sn := ServiceNetwork{
			Network: item.Key.(string),
		}

		if item.Value != nil {
			opts := item.Value.(map[string]interface{})

			if aliases, ok := opts["aliases"]; ok {
				sn.Aliases = toStringSlice(aliases.([]interface{}))
			}
			if ipV4Address, ok := opts["ipv4_address"]; ok {
				sn.IPV4Address = ipV4Address.(string)
			}
			if ipV6Address, ok := opts["ipv6_address"]; ok {
				sn.IPV6Address = ipV6Address.(string)
			}
			if linkLocalIPs, ok := opts["link_local_ips"]; ok {
				sn.LinkLocalIPs = toStringSlice(linkLocalIPs.([]interface{}))
			}
			if priority, ok := opts["priority"]; ok {
				sn.Priority = int64(priority.(uint64))
			}
		}

		nets[i] = sn
	}
	*sn = nets

	return nil
}

func toStringSlice(xs []interface{}) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = x.(string)
	}
	return out
}
