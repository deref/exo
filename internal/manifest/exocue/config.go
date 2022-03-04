package exocue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
)

type Cluster cue.Value
type Stack cue.Value
type Component cue.Value

func Final(v cue.Value) ast.Node {
	return StructToFile(v.Syntax(cue.Final()))
}

func (s Stack) Component(name string) Component {
	v := cue.Value(s)
	componentPath := cue.MakePath(cue.Str("components"), cue.Str(name))
	res := v.LookupPath(componentPath)
	return Component(res)
}
