package compose

import (
	"fmt"
	"sort"
)

type Ulimits []Ulimit

type Limits struct {
	Hard int64
	Soft int64
}

type Ulimit struct {
	Name string
	Limits
}

func (u *Ulimits) UnmarshalYAML(unmarshal func(interface{}) error) error {
	limitMap := map[string]interface{}{}
	if err := unmarshal(&limitMap); err != nil {
		return err
	}

	i, keys := 0, make([]string, len(limitMap))
	for k := range limitMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	ulimits := make([]Ulimit, len(keys))
	for i, k := range keys {
		v := limitMap[k]
		var limits Limits
		switch v := v.(type) {
		case int:
			limits = Limits{Hard: int64(v), Soft: int64(v)}
		case map[string]interface{}:
			hard, okay := v["hard"].(int)
			if !okay {
				return fmt.Errorf("invalid value for hard ulimit: %v", v["hard"])
			}
			soft, okay := v["soft"].(int)
			if !okay {
				return fmt.Errorf("invalid value for soft ulimit: %v", v["soft"])
			}
			limits = Limits{Hard: int64(hard), Soft: int64(soft)}
		default:
			return fmt.Errorf("unexpected value %t", v)
		}

		ulimits[i] = Ulimit{
			Name:   k,
			Limits: limits,
		}
	}
	*u = ulimits
	return nil
}
