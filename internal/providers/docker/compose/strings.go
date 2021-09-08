package compose

import (
	"github.com/goccy/go-yaml"
)

type StringOrStringSlice []string

func (ss StringOrStringSlice) MarshalYAML() (interface{}, error) {
	if len(ss) == 1 {
		return ss[0], nil
	}
	return []string(ss), nil
}

func (ss *StringOrStringSlice) UnmarshalYAML(b []byte) error {
	var asSlice []interface{}
	if err := yaml.Unmarshal(b, &asSlice); err == nil {
		stringSlice := make([]string, len(asSlice))
		for i, s := range asSlice {
			stringSlice[i] = s.(string)
		}
		*ss = StringOrStringSlice(stringSlice)
		return nil
	}

	var asString string
	err := yaml.Unmarshal(b, &asString)
	if err == nil {
		*ss = []string{asString}
	}

	return err
}
