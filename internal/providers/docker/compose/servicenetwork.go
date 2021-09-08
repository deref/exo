package compose

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
)

type serviceNetworkWithoutName struct {
	Aliases      []string `json:"aliases,omitempty"`
	IPV4Address  string   `json:"ipv4_address,omitempty"`
	IPV6Address  string   `json:"ipv6_address,omitempty"`
	LinkLocalIPs []string `json:"link_local_ips,omitempty"`
	Priority     int64    `json:"priority,omitempty"`
}

type ServiceNetwork struct {
	Network string
	serviceNetworkWithoutName
}

type ServiceNetworks []ServiceNetwork

func (sn *ServiceNetworks) UnmarshalYAML(b []byte) error {
	var asStrings []string
	if err := yaml.Unmarshal(b, &asStrings); err == nil {
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
	if err := yaml.Unmarshal(b, &asMap); err != nil {
		return fmt.Errorf("unmarshalling networks: %w", err)
	}

	nets := make([]ServiceNetwork, len(asMap))
	for i, item := range asMap {
		sn := ServiceNetwork{
			Network: item.Key.(string),
		}

		if item.Value != nil {
			opts, ok := item.Value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("could not unmarshal network item %s", sn.Network)
			}

			if jsonBytes, err := json.MarshalIndent(opts, "", "  "); err == nil {
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Println("Error printing opts:", err)
			}

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

func (networks ServiceNetworks) MarshalYAML() (interface{}, error) {
	sns := map[string]serviceNetworkWithoutName{}
	for _, sn := range networks {
		sns[sn.Network] = sn.serviceNetworkWithoutName
	}
	return sns, nil
}

func toStringSlice(xs []interface{}) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = x.(string)
	}
	return out
}
