package compose

type Volume struct {
	Key string `yaml:"-"`

	Driver     String     `yaml:"driver,omitempty"`
	DriverOpts Dictionary `yaml:"driver_opts,omitempty"`
	// TODO: external
	Labels Dictionary `yaml:"labels,omitempty"`
	Name   String     `yaml:"name,omitempty"`
}

func (v *Volume) Interpolate(env Environment) error {
	return interpolateStruct(v, env)
}
