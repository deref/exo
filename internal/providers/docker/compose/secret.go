package compose

type Secret struct {
	File     string `yaml:"file,omitempty"`
	External bool   `yaml:"external,omitempty"`
	Name     string `yaml:"name,omitempty"`
}
