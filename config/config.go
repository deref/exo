package config

// TODO: Add yaml tags to struct fields.

var Version = "0.1"

type Config struct {
	Exo        string
	Components []Component
}

type Component struct {
	Name string
	Type string
	Spec interface{}
}

func NewConfig() *Config {
	return &Config{
		Exo: Version,
	}
}
