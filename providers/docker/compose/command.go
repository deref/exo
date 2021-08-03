package compose

import "mvdan.cc/sh/v3/shell"

type Command []string

func (cmd Command) MarshalYAML() (interface{}, error) {
	return []string(cmd), nil
}

func (cmd *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strs []string
	err := unmarshal(&strs)
	if err == nil {
		*cmd = Command(strs)
		return nil
	}

	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var noEnv func(string) string
	strs, err = shell.Fields(s, noEnv)
	*cmd = Command(strs)
	return err
}
