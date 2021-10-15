package compose

type Healthcheck struct {
	Test        Command  `yaml:"test,omitempty"`
	Interval    Duration `yaml:"interval,omitempty"`
	Timeout     Duration `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod Duration `yaml:"start_period,omitempty"`
}
