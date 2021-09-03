package compose

import "fmt"

type EnvironmentFiles []string

func (ef *EnvironmentFiles) UnmarshalYAML(unmarshal func(interface{}) error) error {
	files := make([]string, 0)
	if err := unmarshal(&files); err == nil {
		*ef = files
		return nil
	}

	var file string
	if err := unmarshal(&file); err != nil {
		return fmt.Errorf("unmarshalling environment files: %w", err)
	}
	*ef = []string{file}
	return nil
}
