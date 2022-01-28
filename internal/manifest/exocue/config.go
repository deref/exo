package exocue

import "cuelang.org/go/cue"

type Stack struct {
	v cue.Value
}

func NewStack(v cue.Value) *Stack {
	return &Stack{
		v: v,
	}
}

func (s *Stack) Eval() cue.Value {
	return s.v.Eval()
}

func (s *Stack) evalPath(selectors ...cue.Selector) cue.Value {
	path := cue.MakePath(selectors...)
	return s.v.LookupPath(path).Eval()
}

func (s *Stack) Component(name string) cue.Value {
	return s.evalPath(cue.Str("components"), cue.Str(name))
}
