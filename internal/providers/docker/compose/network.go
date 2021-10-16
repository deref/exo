package compose

type Network struct {
	Key string `yaml:"-"`

	Name       String     `yaml:"name,omitempty"`
	Driver     String     `yaml:"driver,omitempty"`
	DriverOpts Dictionary `yaml:"driver_opts,omitempty"`
	Attachable Bool       `yaml:"attachable,omitempty"`
	EnableIPv6 Bool       `yaml:"enable_ipv6,omitempty"`
	Internal   Bool       `yaml:"internal,omitempty"`
	Labels     Dictionary `yaml:"labels,omitempty"`
	External   Bool       `yaml:"external,omitempty"`
}
