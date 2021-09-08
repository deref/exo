package compose

import (
	"encoding/json"
	"fmt"

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
		return err
	}

	nets := make([]ServiceNetwork, len(asMap))
	for i, item := range asMap {
		sn := ServiceNetwork{
			Network: item.Key.(string),
		}

		if item.Value != nil {
			opts := item.Value.(map[string]interface{})

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

func (sn ServiceNetworks) MarshalYAML() (interface{}, error) {
	out := make(map[string]interface{}, len(sn))
	for _, n := range sn {
		out[n.Network] = map[string]interface{}{}
		if n.Aliases != nil {
			out["aliases"] = n.Aliases
		}
		if n.IPV4Address != "" {
			out["ipv4_address"] = n.IPV4Address
		}
		if n.IPV6Address != "" {
			out["ipv6_address"] = n.IPV6Address
		}
		if n.LinkLocalIPs != nil {
			out["link_local_ips"] = n.LinkLocalIPs
		}
		if n.Priority != 0 {
			out["priority"] = n.Priority
		}
	}
	return out, nil
}

func toStringSlice(xs []interface{}) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = x.(string)
	}
	return out
}
