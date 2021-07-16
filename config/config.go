package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var Version = "0.1"

type Config struct {
	Exo        string      `hcl:"exo"`
	Components []Component `hcl:"component,block"`
}

type Component struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type,label"`
	Spec string `hcl:"spec"` // TODO: Custom unmarshalling to allow convenient json representation.
}

func NewConfig() *Config {
	return &Config{
		Exo: Version,
	}
}

func Parse(bs []byte) (*Config, error) {
	var cfg Config
	evalCtx := &hcl.EvalContext{
		Functions: map[string]function.Function{
			"jsonencode": stdlib.JSONEncodeFunc,
		},
	}
	if err := hclsimple.Decode("exo.hcl", bs, evalCtx, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
