package exocue

import (
	"encoding/json"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
)

type Configuration cue.Value
type Cluster cue.Value
type Stack cue.Value
type Component cue.Value

func Final(v cue.Value) ast.Node {
	return StructToFile(v.Syntax(cue.Final()))
}

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

func (c Component) Model() json.RawMessage {
	var model json.RawMessage
	if err := lookup(cue.Value(c), cue.Str("model")).Decode(&model); err != nil {
		panic(err)
	}
	return model
}

// TODO: Should be possible to ensure component is valid before calling
// Environment, so returning an error wouldn't be necessary.
func (c Component) Environment() (m map[string]string, err error) {
	err = lookup(cue.Value(c), cue.Str("environment")).Decode(&m)
	return
}
