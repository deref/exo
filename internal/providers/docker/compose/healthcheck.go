package compose

type Healthcheck struct {
	Test        Command  `yaml:"test,omitempty"`
	Interval    Duration `yaml:"interval,omitempty"`
	Timeout     Duration `yaml:"timeout,omitempty"`
	Retries     Int      `yaml:"retries,omitempty"`
	StartPeriod Duration `yaml:"start_period,omitempty"`
}

func (hc *Healthcheck) Interpolate(env Environment) error {
	return interpolateStruct(hc, env)
}
