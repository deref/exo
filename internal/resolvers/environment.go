package resolvers

import (
	"fmt"
	"sort"

	. "github.com/deref/exo/internal/scalars"
)

type EnvironmentResolver struct {
	Variables []*VariableResolver
}

func JSONObjectToEnvironment(obj JSONObject, source string) (*EnvironmentResolver, error) {
	variables := make([]*VariableResolver, 0, len(obj))
	for k, v := range obj {
		vs, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("environment variable %q value is not a string", k)
		}
		variables = append(variables, &VariableResolver{
			Name:   k,
			Value:  vs,
			Source: source,
		})
	}
	sort.Sort(VariablesByName(variables))
	environment := &EnvironmentResolver{
		Variables: variables,
	}
	return environment, nil
}

func EnvMapToEnvironment(m map[string]string, source string) *EnvironmentResolver {
	variables := make([]*VariableResolver, 0, len(m))
	for k, v := range m {
		variables = append(variables, &VariableResolver{
			Name:   k,
			Value:  v,
			Source: source,
		})
	}
	sort.Sort(VariablesByName(variables))
	environment := &EnvironmentResolver{
		Variables: variables,
	}
	return environment
}

func (r *EnvironmentResolver) AsMap() JSONObject {
	obj := make(JSONObject)
	for _, variable := range r.Variables {
		obj[variable.Name] = variable.Value
	}
	return obj
}
