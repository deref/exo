package compose

import (
	"fmt"
	"sort"
	"strings"
)

// string->string mapping that may be encoded as either a map or an array of
// pairs each encoded as "name=value".
type Dictionary map[string]string

func (dict Dictionary) MarshalYAML() (interface{}, error) {
	return map[string]string(dict), nil
}

func (dict *Dictionary) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data interface{}
	if err := unmarshal(&data); err != nil {
		return err
	}

	res := make(map[string]string)
	switch data := data.(type) {
	case map[string]interface{}:
		for k, v := range data {
			var ok bool
			res[k], ok = v.(string)
			if !ok {
				return fmt.Errorf("expected values to be string, got %v", v)
			}
		}
	case []interface{}:
		for _, elem := range data {
			s, ok := elem.(string)
			if !ok {
				return fmt.Errorf("expected elements to be string, got %T", elem)
			}
			parts := strings.SplitN(s, "=", 2)
			k := parts[0]
			v := ""
			if len(parts) == 2 {
				v = parts[1]
			}
			res[k] = v
		}
	default:
		return fmt.Errorf("expected map or array, got %T", data)
	}

	*dict = res
	return nil
}

func (dict Dictionary) Slice() []string {
	m := map[string]string(dict)
	res := make([]string, len(m))
	i := 0
	for k, v := range m {
		res[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	sort.Strings(res)
	return res
}
