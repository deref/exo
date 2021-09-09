package compose

import "sort"

type ServiceNetworkWithoutName struct {
	Aliases      []string `yaml:"aliases,omitempty"`
	IPV4Address  string   `yaml:"ipv4_address,omitempty"`
	IPV6Address  string   `yaml:"ipv6_address,omitempty"`
	LinkLocalIPs []string `yaml:"link_local_ips,omitempty"`
	Priority     int64    `yaml:"priority,omitempty"`
}

type ServiceNetwork struct {
	Network string
	ServiceNetworkWithoutName
}

type ServiceNetworks []ServiceNetwork

func (sns ServiceNetworks) MarshalYAML() (interface{}, error) {
	m := map[string]ServiceNetworkWithoutName{}
	for _, sn := range sns {
		m[sn.Network] = sn.ServiceNetworkWithoutName
	}
	return m, nil
}

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

	var asMap map[string]*ServiceNetworkWithoutName
	if err := unmarshal(&asMap); err != nil {
		return err
	}

	i, keys := 0, make([]string, len(asMap))
	for k := range asMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	nets := []ServiceNetwork{}
	for _, key := range keys {
		item := asMap[key]
		if item == nil {
			item = &ServiceNetworkWithoutName{}
		}
		nets = append(nets, ServiceNetwork{
			Network:                   key,
			ServiceNetworkWithoutName: *item,
		})
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
