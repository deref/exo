package compose

import "gopkg.in/yaml.v3"

type String struct {
	Style      yaml.Style
	Expression string
	Value      string
}

func MakeString(s string) String {
	return String{
		Expression: s,
		Value:      s,
	}
}

func (s *String) UnmarshalYAML(node *yaml.Node) error {
	s.Style = node.Style
	err := node.Decode(&s.Expression)
	s.Value = s.Expression
	return err
}

func (s String) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Style: s.Style,
		Value: s.Expression,
	}, nil
}
