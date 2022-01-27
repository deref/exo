package exocue

import "cuelang.org/go/cue"

type Configuration struct {
	v cue.Value
}

func NewConfiguration(v cue.Value) *Configuration {
	return &Configuration{
		v: v,
	}
}

func (cfg *Configuration) Eval() cue.Value {
	return cfg.v.Eval()
}

func (cfg *Configuration) evalPath(selectors ...cue.Selector) cue.Value {
	path := cue.MakePath(selectors...)
	return cfg.v.LookupPath(path).Eval()
}

func (cfg *Configuration) Component(name string) cue.Value {
	return cfg.evalPath(cue.Str("$stack"), cue.Str("components"), cue.Str(name))
}

func (cfg *Configuration) ComponentSpec(name string) cue.Value {
	return cfg.evalPath(cue.Str("$stack"), cue.Str("components"), cue.Str(name), cue.Str("spec"))
}
