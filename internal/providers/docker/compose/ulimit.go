package compose

import (
	"github.com/goccy/go-yaml"
)

type Ulimits []Ulimit

type Ulimit struct {
	Name string
	Hard int64
	Soft int64
}

func (u *Ulimits) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asMap yaml.MapSlice
	if err := unmarshal(&asMap); err != nil {
		return err
	}

	ulimits := make([]Ulimit, len(asMap))
	for i, item := range asMap {
		ulimit := Ulimit{
			Name: item.Key.(string),
		}
		switch val := item.Value.(type) {
		case uint64:
			ulimit.Hard = int64(val)
			ulimit.Soft = int64(val)
		case map[string]interface{}:
			ulimit.Hard = int64(val["hard"].(uint64))
			ulimit.Soft = int64(val["soft"].(uint64))
		}
		ulimits[i] = ulimit
	}
	*u = ulimits

	return nil
}
