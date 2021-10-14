package compose

import "gopkg.in/yaml.v3"

type Strings struct {
	IsSequence bool
	Values     []string
}

func (ss Strings) MarshalYAML() (interface{}, error) {
	if ss.IsSequence || len(ss.Values) != 1 {
		return ss.Values, nil
	}
	return ss.Values[0], nil
}

func (ss *Strings) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if node.Decode(&s) == nil {
		ss.Values = []string{s}
		return nil
	}
	ss.IsSequence = true
	return node.Decode(&ss.Values)
}
