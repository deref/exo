package compose

import "gopkg.in/yaml.v3"

type Command struct {
	// IsShellForm is true if the command was provided as a string rather than an
	// array of strings. This indicates that it should be passed to the image's
	// shell as the first argument.
	IsShellForm bool
	Parts       Strings
}

func (cmd Command) MarshalYAML() (interface{}, error) {
	if cmd.IsShellForm {
		return cmd.Parts[0], nil
	}
	// In parsed form, prefer flow-style with double quotes to
	// match typical Dockerfile style.
	// TODO: Preserve style of source.
	content := make([]*yaml.Node, len(cmd.Parts))
	for i, part := range cmd.Parts {
		partNode := &yaml.Node{}
		if err := partNode.Encode(part); err != nil {
			panic(err)
		}
		partNode.Style = yaml.DoubleQuotedStyle
		content[i] = partNode
	}
	node := &yaml.Node{
		Kind:    yaml.SequenceNode,
		Style:   yaml.FlowStyle,
		Content: content,
	}
	return node, nil
}

func (cmd *Command) UnmarshalYAML(node *yaml.Node) error {
	var strs []String
	err := node.Decode(&strs)
	if err == nil {
		for i := range strs {
			strs[i].Style = 0
		}
		cmd.Parts = strs
		return nil
	}

	var s String
	if err := node.Decode(&s); err != nil {
		return err
	}
	s.Style = 0
	cmd.IsShellForm = true
	cmd.Parts = []String{s}
	return nil
}

func (cmd *Command) Interpolate(env Environment) error {
	for i := range cmd.Parts {
		if err := cmd.Parts[i].Interpolate(env); err != nil {
			return err
		}
	}
	return nil
}
