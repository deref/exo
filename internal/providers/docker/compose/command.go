package compose

import (
	"github.com/mattn/go-shellwords"
)

// FIXME: this incorrectly treats the "command" passed to docker as a meaningful
// shell command that can be parsed into "shell words". In reality the command
// is an arbitrary string that is passed to the entrypoint of the container as
// an argument. It could be a shell command but it could also be a python
// program or a Dan Brown novel.
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

	strs, err = shellwords.Parse(s)
	*cmd = Command(strs)
	return err
}
