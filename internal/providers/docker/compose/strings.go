package compose

type StringOrStringSlice []string

func (ss StringOrStringSlice) MarshalYAML() (interface{}, error) {
	if len(ss) == 1 {
		return ss[0], nil
	}
	return []string(ss), nil
}

func (ss *StringOrStringSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asSlice []interface{}
	if err := unmarshal(&asSlice); err == nil {
		stringSlice := make([]string, len(asSlice))
		for i, s := range asSlice {
			stringSlice[i] = s.(string)
		}
		*ss = StringOrStringSlice(stringSlice)
		return nil
	}

	var asString string
	err := unmarshal(&asString)
	if err == nil {
		*ss = []string{asString}
	}

	return err
}
