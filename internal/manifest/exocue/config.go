package exocue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
)

func Final(v cue.Value) ast.Node {
	return StructToFile(v.Syntax(cue.Final()))
}

type Stack cue.Value

func (s Stack) Final() ast.Node {
	return Final(cue.Value(s))
}

func (s Stack) Component(name string) Component {
	v := cue.Value(s)
	componentPath := cue.MakePath(cue.Str("components"), cue.Str(name))
	res := v.LookupPath(componentPath)
	return Component(res)
}

type Component cue.Value

func (c Component) Final() ast.Node {
	return Final(cue.Value(c))
}
