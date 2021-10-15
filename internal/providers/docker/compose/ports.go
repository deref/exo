package compose

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type PortMappings []PortMapping

type PortMapping struct {
	IsShortForm bool
	PortMappingLongForm
}

type PortMappingLongForm struct {
	Target    PortRange `yaml:"target,omitempty"`
	Published PortRange `yaml:"published,omitempty"`
	HostIP    string    `yaml:"host_ip,omitempty"`
	Protocol  string    `yaml:"protocol,omitempty"`
	Mode      string    `yaml:"mode,omitempty"`
}

func ParsePortMappings(short string) (mappings PortMappings, err error) {
	parts := strings.Split(short, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		mapping, err := ParsePortMapping(part)
		if err != nil {
			return PortMappings{}, err
		}
		mappings = append(mappings, PortMapping{
			IsShortForm:         true,
			PortMappingLongForm: mapping,
		})
	}
	return mappings, nil
}

func ParsePortMapping(short string) (PortMappingLongForm, error) {
	submatches := portRegexp.FindStringSubmatch(short)
	if len(submatches) == 0 {
		return PortMappingLongForm{}, errors.New("invalid port mapping syntax")
	}

	result := make(map[string]string)
	for i, name := range portRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = submatches[i]
		}
	}

	var mapping PortMappingLongForm
	var err error
	mapping.HostIP = result["ip"]
	mapping.Published, err = ParsePortRange(result["published"])
	if err != nil {
		return PortMappingLongForm{}, fmt.Errorf("invalid published port range: %w", err)
	}
	mapping.Target, err = ParsePortRange(result["target"])
	if err != nil {
		return PortMappingLongForm{}, fmt.Errorf("invalid target port range: %w", err)
	}
	mapping.Protocol = result["protocol"]

	if mapping.HostIP != "" && net.ParseIP(mapping.HostIP) == nil {
		return PortMappingLongForm{}, fmt.Errorf("invalid IP: %s", mapping.HostIP)
	}

	return mapping, nil
}

// https://regex101.com/r/qvbqTT/2
var portRegexp = regexp.MustCompile(`^((?P<ip>[a-fA-F\d.:]+?):)??((?P<published>([-\d]+)?):)?(?P<target>[-\d]+)(/(?P<protocol>.+))?$`)

func (mappings PortMappings) MarshalYAML() (interface{}, error) {
	res := make([]interface{}, len(mappings))
	for i, x := range mappings {
		res[i] = x // TODO: Marshal to short syntax if possible.
	}
	return res, nil
}

func (pm *PortMapping) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if node.Decode(&s) == nil {
		pm.IsShortForm = true
		var err error
		pm.PortMappingLongForm, err = ParsePortMapping(s)
		return err
	}
	return node.Decode(&pm.PortMappingLongForm)
}

func (pm PortMapping) MarshalYAML() (interface{}, error) {
	if pm.IsShortForm {
		return pm.Target, nil
	}
	return pm.PortMappingLongForm, nil
}

type PortRange struct {
	Min uint16
	Max uint16
}

func ParsePortRange(numbers string) (res PortRange, err error) {
	if numbers == "" {
		return PortRange{}, nil
	}
	parts := strings.SplitN(numbers, "-", 2)
	if len(parts) == 1 {
		parts = append(parts, parts[0])
	}
	for i, dest := range []*uint16{&res.Min, &res.Max} {
		var n int
		n, err = strconv.Atoi(parts[i])
		*dest = uint16(n)
		if err != nil {
			return
		}
	}
	return
}

func FormatPort(num uint16, protocol string) string {
	res := strconv.Itoa(int(num))
	if protocol != "" {
		res += "/" + protocol
	}
	return res
}

func (rng *PortRange) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}
	var err error
	*rng, err = ParsePortRange(s)
	return err
}

func (rng PortRange) MarshalYAML() (interface{}, error) {
	if rng.Min == rng.Max {
		return rng.Min, nil
	}
	return fmt.Sprintf("%d-%d", rng.Min, rng.Max), nil
}
