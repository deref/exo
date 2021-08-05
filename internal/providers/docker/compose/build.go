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
	Context    string     `yaml:"context"`
	Dockerfile string     `yaml:"dockerfile"`
	Args       Dictionary `yaml:"args"`
	CacheFrom  []string   `yaml:"cache_from"`
	ExtraHosts []string   `yaml:"extra_hosts"`
	Isolation  string     `yaml:"isolation"`
	Labels     Dictionary `yaml:"labels"`
	ShmSize    Bytes      `yaml:"shm_size"`
	Target     string     `yaml:"target"`
}
