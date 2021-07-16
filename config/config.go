package config

var Version = "0.1"

type Config struct {
	Exo        string      `hcl:"exo"`
	Components []Component `hcl:"component,block"`
}

type Component struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type,label"`
	Spec string `hcl:"spec"` // TODO: Custom unmarshalling to allow convenient json representation.
}

func NewConfig() *Config {
	return &Config{
		Exo: Version,
	}
}
