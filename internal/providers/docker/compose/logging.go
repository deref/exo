package compose

type Logging struct {
	Driver  String     `yaml:"driver,omitempty"`
	Options Dictionary `yaml:"options,omitempty"`
}

func (l *Logging) Interpolate(env Environment) error {
	return interpolateStruct(l, env)
}
