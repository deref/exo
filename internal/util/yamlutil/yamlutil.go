package yamlutil

import "github.com/goccy/go-yaml"

func MarshalString(v interface{}) (string, error) {
	bs, err := yaml.Marshal(v)
	return string(bs), err
}

func MustMarshalString(v interface{}) string {
	s, err := MarshalString(v)
	if err != nil {
		panic(err)
	}
	return s
}

func UnmarshalString(s string, v interface{}) error {
	return yaml.Unmarshal([]byte(s), v)
}

func MustUnmarshalString(s string, v interface{}) {
	if err := UnmarshalString(s, v); err != nil {
		panic(err)
	}
}
