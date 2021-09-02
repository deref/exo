package compose

type Command struct {
	Parts []string
	// IsShellForm is true if the command was provided as a string rather than an
	// array of strings. This indicates that it should be passed to the image's
	// shell as the first argument.
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
