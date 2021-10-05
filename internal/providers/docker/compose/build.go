package compose

type BuildTemplate struct {
	IsShortForm bool
	BuildTemplateLongForm
}

type BuildTemplateLongForm struct {
	Context    string     `yaml:"context,omitempty"`
	Dockerfile string     `yaml:"dockerfile,omitempty"`
	Args       Dictionary `yaml:"args,omitempty"`
	CacheFrom  []string   `yaml:"cache_from,omitempty"`
	ExtraHosts []string   `yaml:"extra_hosts,omitempty"`
	Isolation  string     `yaml:"isolation,omitempty"`
	Labels     Dictionary `yaml:"labels,omitempty"`
	ShmSize    string     `yaml:"shm_size,omitempty"`
	Target     string     `yaml:"target,omitempty"`
}

func (b *BuildTemplate) MarshalYAML() (interface{}, error) {
	return b.BuildTemplateLongForm, nil
}

func (b *BuildTemplate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err == nil {
		b.IsShortForm = true
		b.Context = s
		return nil
	}
	return unmarshal(&b.BuildTemplateLongForm)
}

type Build struct {
	Context    string     `yaml:"context,omitempty"`
	Dockerfile string     `yaml:"dockerfile,omitempty"`
	Args       Dictionary `yaml:"args,omitempty"`
	CacheFrom  []string   `yaml:"cache_from,omitempty"`
	ExtraHosts []string   `yaml:"extra_hosts,omitempty"`
	Isolation  string     `yaml:"isolation,omitempty"`
	Labels     Dictionary `yaml:"labels,omitempty"`
	ShmSize    Bytes      `yaml:"shm_size,omitempty"`
	Target     string     `yaml:"target,omitempty"`
}
