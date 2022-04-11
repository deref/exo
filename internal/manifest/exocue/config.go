package exocue

import (
	"cuelang.org/go/cue"
	. "github.com/deref/exo/internal/scalars"
)

type Configuration cue.Value
type Cluster cue.Value
type Stack cue.Value
type Component cue.Value

func lookup(v cue.Value, path ...cue.Selector) cue.Value {
	return v.LookupPath(cue.MakePath(path...))
}

func (cfg Configuration) Cluster() Cluster {
	return Cluster(lookup(cue.Value(cfg), cue.Str("$cluster")))
}

func (cfg Configuration) Stack() Stack {
	return Stack(lookup(cue.Value(cfg), cue.Str("$stack")))
}

func (cfg Configuration) Component(id string) Component {
	return Component(lookup(cue.Value(cfg), cue.Str("$components"), cue.Str(id)))
}

func (c Component) Model() RawJSON {
	var model RawJSON
	if err := lookup(cue.Value(c), cue.Str("model")).Decode(&model); err != nil {
		panic(err)
	}
	return model
}

// TODO: Should be possible to ensure component is valid before calling
// Environment, so returning an error wouldn't be necessary.
func (s Stack) FullEnvironment() (m map[string]string, err error) {
	err = lookup(cue.Value(s), cue.Str("fullEnvironment")).Decode(&m)
	return
}
func (c Component) FullEnvironment() (m map[string]string, err error) {
	err = lookup(cue.Value(c), cue.Str("fullEnvironment")).Decode(&m)
	return
}
