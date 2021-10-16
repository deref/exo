package compose

import "gopkg.in/yaml.v3"

type Build struct {
	ShortForm String
	BuildLongForm
}

type BuildLongForm struct {
	Context    String     `yaml:"context,omitempty"`
	Dockerfile String     `yaml:"dockerfile,omitempty"`
	Args       Dictionary `yaml:"args,omitempty"`
	CacheFrom  Strings    `yaml:"cache_from,omitempty"`
	ExtraHosts Strings    `yaml:"extra_hosts,omitempty"`
	Isolation  String     `yaml:"isolation,omitempty"`
	Labels     Dictionary `yaml:"labels,omitempty"`
	ShmSize    Bytes      `yaml:"shm_size,omitempty"`
	Target     String     `yaml:"target,omitempty"`
}

func (b Build) MarshalYAML() (interface{}, error) {
	if b.ShortForm.Expression != "" {
		return b.ShortForm.Expression, nil
	}
	return b.BuildLongForm, nil
}

func (b *Build) UnmarshalYAML(node *yaml.Node) error {
	var err error
	if node.Tag == "!!str" {
		err = node.Decode(&b.ShortForm)
	} else {
		err = node.Decode(&b.BuildLongForm)
	}
	_ = b.Interpolate(ErrEnvironment)
	return err
}

func (b *Build) Interpolate(env Environment) error {
	if b.ShortForm.Expression != "" {
		err := b.ShortForm.Interpolate(env)
		b.Context = b.ShortForm
		return err
	} else {
		return b.BuildLongForm.Interpolate(env)
	}
}

func (b *BuildLongForm) Interpolate(env Environment) error {
	return interpolateStruct(b, env)
}
