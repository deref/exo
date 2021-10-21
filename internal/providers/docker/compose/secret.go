package compose

type Secret struct {
	Key string `yaml:"-"`

	File     String `yaml:"file,omitempty"`
	External Bool   `yaml:"external,omitempty"`
	Name     String `yaml:"name,omitempty"`
}

func (s *Secret) Interpolate(env Environment) error {
	return interpolateStruct(s, env)
}
