package compose

type Volume struct {
	Driver     string     `yaml:"driver,omitempty"`
	DriverOpts Dictionary `yaml:"driver_opts,omitempty"`
	// TODO: external
	Labels Dictionary `yaml:"labels,omitempty"`
	Name   string     `yaml:"name,omitempty"`
}
