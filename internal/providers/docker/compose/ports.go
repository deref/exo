package compose

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type PortMappings []PortMapping

type PortMapping struct {
	Target    string `yaml:"target"`
	Published string `yaml:"published"`
	HostIP    string `yaml:"host_ip"`
	Protocol  string `yaml:"protocol"`
	Mode      string `yaml:"mode"`
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
			return nil, err
		}
		mappings = append(mappings, mapping)
	}
	return mappings, nil
}

func ParsePortMapping(short string) (PortMapping, error) {
	submatches := portRegexp.FindStringSubmatch(short)
	if len(submatches) == 0 {
		return PortMapping{}, errors.New("invalid port mapping syntax")
	}

	result := make(map[string]string)
	for i, name := range portRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = submatches[i]
		}
	}

	var mapping PortMapping
	mapping.HostIP = result["ip"]
	mapping.Published = result["published"]
	mapping.Target = result["target"]
	mapping.Protocol = result["protocol"]
	return mapping, nil
}

// https://regex101.com/r/qvbqTT/1
var portRegexp = regexp.MustCompile(`^((?P<ip>(\d+\.\d+\.\d+\.\d+)|(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))):)?((?P<published>([-\d]+)?):)?(?P<target>[-\d]+)(/(?P<protocol>.+))?$`)

func (mappings PortMappings) MarshalYAML() (interface{}, error) {
	res := make([]interface{}, len(mappings))
	for i, x := range mappings {
		res[i] = x // TODO: Marshal to short syntax if possible.
	}
	return res, nil
}

func (mappings *PortMappings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*mappings = []PortMapping{}
	var s string
	err := unmarshal(&s)
	if err == nil {
		*mappings, err = ParsePortMappings(s)
	} else {
		var xs []portMappingMarshaller
		err = unmarshal(&xs)
		for _, x := range xs {
			*mappings = append(*mappings, PortMapping(x))
		}
	}
	return err
}

type portMappingMarshaller PortMapping

func (marshaller *portMappingMarshaller) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	var mapping PortMapping
	err := unmarshal(&s)
	if err == nil {
		mapping, err = ParsePortMapping(s)
	} else {
		err = unmarshal(&mapping)
	}
	*marshaller = portMappingMarshaller(mapping)
	return err
}

func ParsePortRange(numbers string, protocol string) (res PortRange, err error) {
	res.Protocol = protocol
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

type PortRange struct {
	Min      uint16
	Max      uint16
	Protocol string
}

func FormatPort(num uint16, protocol string) string {
	res := strconv.Itoa(int(num))
	if protocol != "" {
		res += "/" + protocol
	}
	return res
}
