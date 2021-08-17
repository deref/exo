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
		cfg.Context.Expression = s
	} else if err := unmarshal(&cfg); err != nil {
		return nil
	}
	*dict = Build(cfg)
	return nil
}

type BuildConfig struct {
	Context    String     `yaml:"context"`
	Dockerfile String     `yaml:"dockerfile"`
	Args       Dictionary `yaml:"args"`
	CacheFrom  []String   `yaml:"cache_from"`
	ExtraHosts []String   `yaml:"extra_hosts"`
	Isolation  String     `yaml:"isolation"`
	Labels     Dictionary `yaml:"labels"`
	ShmSize    Bytes      `yaml:"shm_size"`
	Target     String     `yaml:"target"`
}
