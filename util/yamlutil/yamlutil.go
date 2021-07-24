package yamlutil

import "github.com/goccy/go-yaml"

func MustMarshalString(v interface{}) string {
	bs, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
