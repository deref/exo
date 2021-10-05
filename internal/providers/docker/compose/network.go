package compose

import "github.com/goccy/go-yaml"

type NetworkTemplate struct {
	Name     string        `yaml:"name,omitempty"`
	Driver   string        `yaml:"driver,omitempty"`
	MapSlice yaml.MapSlice `yaml:",inline"`
}

type Network struct {
	// Name is the actual name of the docker network. The docker-compose network name, which can
	// be referenced by individual services, is the component name.
	Name       string            `yaml:"name,omitempty"`
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	Attachable bool              `yaml:"attachable,omitempty"`
	EnableIPv6 bool              `yaml:"enable_ipv6,omitempty"`
	Internal   bool              `yaml:"internal,omitempty"`
	Labels     Dictionary        `yaml:"labels,omitempty"`
	External   bool              `yaml:"external,omitempty"`
}
