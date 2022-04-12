package resolvers

import (
	"context"
	"fmt"
	"strings"

	. "github.com/deref/exo/internal/scalars"
	"golang.org/x/exp/slices"
)

type EnvironmentSource interface {
	Environment(ctx context.Context) (*EnvironmentResolver, error)
}

type EnvironmentResolver struct {
	Parent *EnvironmentResolver
	Source EnvironmentSource
	Locals []*EnvironmentVariableResolver
}

type EnvironmentVariableResolver struct {
	Name   string
	Value  *string
	Source EnvironmentSource
}

func sortEnvironmentVariables(variables []*EnvironmentVariableResolver) {
	slices.SortFunc(variables, func(a, b *EnvironmentVariableResolver) bool {
		return strings.Compare(a.Name, b.Name) < 0
	})
}

func (r *EnvironmentResolver) initLocalsFromJSONObject(obj JSONObject) {
	r.Locals = make([]*EnvironmentVariableResolver, 0, len(obj))
	for k, v := range obj {
		local := &EnvironmentVariableResolver{
			Name:   k,
			Source: r.Source,
		}
		switch v := v.(type) {
		case nil:
			local.Value = nil
		case string:
			local.Value = &v
		default:
			panic(fmt.Errorf("variable %q has invalid value type: %T", v))
		}
		r.Locals = append(r.Locals, local)
	}
	sortEnvironmentVariables(r.Locals)
}

func (r *EnvironmentResolver) Variables() []*EnvironmentVariableResolver {
	m := r.variablesMap()
	res := make([]*EnvironmentVariableResolver, len(m))
	for _, v := range m {
		if v.Value == nil {
			continue
		}
		res = append(res, v)
	}
	sortEnvironmentVariables(res)
	return res
}

func (r *EnvironmentResolver) variablesMap() map[string]*EnvironmentVariableResolver {
	variables := make(map[string]*EnvironmentVariableResolver)
	r.gatherVariables(variables)
	return variables
}

func (r *EnvironmentResolver) gatherVariables(variables map[string]*EnvironmentVariableResolver) {
	if r.Parent != nil {
		r.Parent.gatherVariables(variables)
	}
	for _, local := range r.Locals {
		variables[local.Name] = local
	}
}

func (r *EnvironmentResolver) AsMap() JSONObject {
	obj := make(JSONObject)
	for _, variable := range r.variablesMap() {
		obj[variable.Name] = variable.Value
	}
	return obj
}

func (r *EnvironmentResolver) asMap() map[string]string {
	obj := make(map[string]string)
	for _, variable := range r.variablesMap() {
		obj[variable.Name] = *variable.Value
	}
	return obj
}
