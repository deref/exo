package compose

import "gopkg.in/yaml.v3"

type Build struct {
	IsShortForm bool
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
	if b.IsShortForm {
		return b.Context, nil
	}
	return b.BuildLongForm, nil
}

func (b *Build) UnmarshalYAML(node *yaml.Node) error {
	var short String
	err := node.Decode(&short)
	if err == nil {
		b.IsShortForm = true
		b.Context = short
		return nil
	}
	return node.Decode(&b.BuildLongForm)
}
