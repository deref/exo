package exocue

import "cuelang.org/go/cue"

type Stack struct {
	v cue.Value
}

func NewStack(v cue.Value) *Stack {
	dumpValue(v)
	return &Stack{
		v: v,
	}
}

func (s *Stack) Eval() cue.Value {
	return s.v.Eval()
}

func (s *Stack) Component(name string) cue.Value {
	componentPath := cue.MakePath(cue.Str("components"), cue.Str(name))
	res := s.v.LookupPath(componentPath)
	return res
}
