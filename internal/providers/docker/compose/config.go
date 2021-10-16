package compose

type Config struct {
	Key string `yaml:"-"`

	File     String `yaml:"file,omitempty"`
	External Bool   `yaml:"external,omitempty"`
	Name     String `yaml:"name,omitempty"`
}
