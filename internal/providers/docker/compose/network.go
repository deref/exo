package compose

type Network struct {
	Key string `yaml:"-"`

	Name       string     `yaml:"name,omitempty"`
	Driver     string     `yaml:"driver,omitempty"`
	DriverOpts Dictionary `yaml:"driver_opts,omitempty"`
	Attachable bool       `yaml:"attachable,omitempty"`
	EnableIPv6 bool       `yaml:"enable_ipv6,omitempty"`
	Internal   bool       `yaml:"internal,omitempty"`
	Labels     Dictionary `yaml:"labels,omitempty"`
	External   bool       `yaml:"external,omitempty"`
}
