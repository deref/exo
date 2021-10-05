package compose

import "github.com/goccy/go-yaml"

type VolumeTemplate struct {
	Name     string        `yaml:"name,omitempty"`
	MapSlice yaml.MapSlice `yaml:",inline"`
}

type Volume struct {
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	// TODO: external
	Labels Dictionary `yaml:"labels,omitempty"`
	Name   string     `yaml:"name,omitempty"`
}
