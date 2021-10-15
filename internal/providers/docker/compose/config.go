package compose

type Config struct {
	Key string `yaml:"-"`

	File     string `yaml:"file,omitempty"`
	External bool   `yaml:"external,omitempty"`
	Name     string `yaml:"name,omitempty"`
}
