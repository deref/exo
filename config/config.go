package config

// TODO: HCL.

var Version = "0.1"

type Config struct {
	Exo        string
	Components []Component
}

type Component struct {
	Name string
	Type string
	Spec string // TODO: Custom unmarshalling to allow convenient json representation.
}

func NewConfig() *Config {
	return &Config{
		Exo: Version,
	}
}
