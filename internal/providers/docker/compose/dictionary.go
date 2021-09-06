package compose

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

// string->string mapping that may be encoded as either a map or an array of
// pairs each encoded as "name=value". If the equal sign is not supplied,
// the value is treated as nil.
type Dictionary map[string]*string

func (dict Dictionary) MarshalYAML() (interface{}, error) {
	return map[string]*string(dict), nil
}

type stringOrNum string

func (sn *stringOrNum) UnmarshalYAML(b []byte) error {
	// This is necessary as otherwise you end up parsing the number and converting
	// it back to a string. That means a value like "99999999999999.999" becomes
	// "1e+14".
	if _, err := strconv.ParseFloat(string(b), 64); err == nil {
		*sn = stringOrNum(string(b))
		return nil
	}

	var s string
	err := yaml.Unmarshal(b, &s)
	if err == nil {
		*sn = stringOrNum(s)
		return nil
	}
	return fmt.Errorf("unmarshaling stringOrNum: %w", err)
}

func (dict *Dictionary) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data interface{}
	if err := unmarshal(&data); err != nil {
		return err
	}

	res := make(map[string]*string)
	switch data := data.(type) {
	case map[string]interface{}:
		stringOrNumMap := make(map[string]stringOrNum)
		if err := unmarshal(&stringOrNumMap); err != nil {
			return err
		}
		for k, v := range stringOrNumMap {
			s := string(v)
			res[k] = &s
		}
	case []interface{}:
		for _, elem := range data {
			s, ok := elem.(string)
			if !ok {
				return fmt.Errorf("expected elements to be string, got %T", elem)
			}
			parts := strings.SplitN(s, "=", 2)
			k := parts[0]
			switch len(parts) {
			case 1:
				res[k] = nil
			case 2:
				res[k] = &parts[1]
			default:
				panic("unreachable")
			}
		}
	default:
		return fmt.Errorf("expected map or array, got %T", data)
	}

	*dict = res
	return nil
}

func (dict Dictionary) Slice() []string {
	m := map[string]*string(dict)
	res := make([]string, len(m))
	i := 0
	for k, v := range m {
		if v == nil {
			res[i] = k
		} else {
			res[i] = fmt.Sprintf("%s=%s", k, *v)
		}
		i++
	}
	sort.Strings(res)
	return res
}

func (dict Dictionary) WithoutNils() map[string]string {
	m := make(map[string]string, len(dict))
	for k, v := range dict {
		if v == nil {
			continue
		}
		m[k] = *v
	}
	return m
}
