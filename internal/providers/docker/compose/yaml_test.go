package compose

import "github.com/deref/exo/internal/util/yamlutil"

func mustLoadYaml(s string, v interface{}, env Environment) {
	yamlutil.MustUnmarshalString(s, v)
	if err := Interpolate(v, env); err != nil {
		panic(err)
	}
}
