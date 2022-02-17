package resolvers

import "strings"

type VariableResolver struct {
	Name   string
	Value  string
	Source string
}

type VariablesByName []*VariableResolver

func (vars VariablesByName) Len() int      { return len(vars) }
func (vars VariablesByName) Swap(i, j int) { vars[i], vars[j] = vars[j], vars[i] }

func (vars VariablesByName) Less(i, j int) bool {
	return strings.Compare(vars[i].Name, vars[j].Name) < 0
}
