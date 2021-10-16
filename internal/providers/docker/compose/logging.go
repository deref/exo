package compose

type Logging struct {
	Driver  String     `yaml:"driver,omitempty"`
	Options Dictionary `yaml:"options,omitempty"`
}
