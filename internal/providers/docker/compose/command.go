package compose

// FIXME: this incorrectly treats the "command" passed to docker as a meaningful
// shell command that can be parsed into "shell words". In reality the command
// is an arbitrary string that is passed to the entrypoint of the container as
// an argument. It could be a shell command but it could also be a python
// program or a Dan Brown novel.
type Command struct {
	Parts       []string
	IsShellForm bool
}

func (cmd Command) MarshalYAML() (interface{}, error) {
	if cmd.IsShellForm {
		return cmd.Parts[0], nil
	}
	return []string(cmd.Parts), nil
}

func (cmd *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strs []string
	if err := unmarshal(&strs); err == nil {
		*cmd = Command{
			Parts:       strs,
			IsShellForm: false,
		}
		return nil
	}

	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	*cmd = Command{
		Parts:       []string{s},
		IsShellForm: true,
	}
	return nil
}
