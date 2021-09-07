package compose

type Build BuildConfig

func (b Build) MarshalYAML() (interface{}, error) {
	return BuildConfig(b), nil
}

func (dict *Build) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	var cfg BuildConfig
	err := unmarshal(&s)
	if err == nil {
		cfg.Context = s
	} else if err := unmarshal(&cfg); err != nil {
		return nil
	}
	*dict = Build(cfg)
	return nil
}

type BuildConfig struct {
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
