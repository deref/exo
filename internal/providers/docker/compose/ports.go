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
	String
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
	mapping.Published.Min, mapping.Published.Max, err = ParsePortRange(result["published"])
	if err != nil {
		return PortMappingLongForm{}, fmt.Errorf("invalid published port range: %w", err)
	}
	mapping.Target.Min, mapping.Target.Max, err = ParsePortRange(result["target"])
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

func (mappings PortMappings) MarshalYAML() (any, error) {
	res := make([]any, len(mappings))
	for i, x := range mappings {
		res[i] = x // TODO: Marshal to short syntax if possible.
	}
	return res, nil
}

func (pm *PortMappings) Interpolate(env Environment) error {
	return interpolateSlice(*pm, env)
}

func (pm *PortMapping) UnmarshalYAML(node *yaml.Node) error {
	switch node.Tag {
	case "!!int", "!!str":
		pm.IsShortForm = true
		if err := node.Decode(&pm.String); err != nil {
			return err
		}
		_ = pm.Interpolate(ErrEnvironment)
		return nil
	default:
		return node.Decode(&pm.PortMappingLongForm)
	}
}

func (pm *PortMapping) Interpolate(env Environment) error {
	if err := pm.String.Interpolate(env); err != nil {
		return err
	}
	var err error
	pm.PortMappingLongForm, err = ParsePortMapping(pm.String.Value)
	return err
}

func (pm PortMapping) MarshalYAML() (any, error) {
	if pm.IsShortForm {
		return pm.String, nil
	}
	return pm.PortMappingLongForm, nil
}

type PortRange struct {
	String
	Min uint16
	Max uint16
}

type PortRangeWithProtocol struct {
	String
	Min      uint16
	Max      uint16
	Protocol string
}

func ParsePortRange(numbers string) (min uint16, max uint16, err error) {
	if numbers == "" {
		return 0, 0, nil
	}
	parts := strings.SplitN(numbers, "-", 2)
	if len(parts) == 1 {
		parts = append(parts, parts[0])
	}
	for i, dest := range []*uint16{&min, &max} {
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
	if err := node.Decode(&rng.String); err != nil {
		return err
	}
	_ = rng.Interpolate(ErrEnvironment)
	return nil
}

func (rng *PortRange) Interpolate(env Environment) error {
	if err := rng.String.Interpolate(env); err != nil {
		return err
	}
	var err error
	rng.Min, rng.Max, err = ParsePortRange(rng.String.Value)
	return err
}

func (rng PortRange) MarshalYAML() (any, error) {
	if rng.Min == rng.Max {
		return rng.Min, nil
	}
	return fmt.Sprintf("%d-%d", rng.Min, rng.Max), nil
}

func (rng *PortRangeWithProtocol) UnmarshalYAML(node *yaml.Node) error {
	if err := node.Decode(&rng.String); err != nil {
		return err
	}
	_ = rng.Interpolate(ErrEnvironment)
	return nil
}

func (rng *PortRangeWithProtocol) Interpolate(env Environment) error {
	if err := rng.String.Interpolate(env); err != nil {
		return err
	}
	parts := strings.SplitN(rng.String.Value, "/", 2)
	if len(parts) > 1 {
		rng.Protocol = parts[1]
	}
	var err error
	rng.Min, rng.Max, err = ParsePortRange(parts[0])
	return err
}

func (rng PortRangeWithProtocol) MarshalYAML() (any, error) {
	if rng.Min == rng.Max {
		return rng.Min, nil
	}
	s := fmt.Sprintf("%d-%d", rng.Min, rng.Max)
	if rng.Protocol != "" {
		s += "/" + rng.Protocol
	}
	return s, nil
}
