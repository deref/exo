package exocue

import (
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

func (c Component) Environment() (m map[string]string, err error) {
	err = lookup(cue.Value(c), cue.Str("environment")).Decode(&m)
	return
}
