package compose

import "gopkg.in/yaml.v3"

type Strings struct {
	IsSequence bool
	Items      []string
}

func MakeStrings(items ...string) Strings {
	res := Strings{
		Items: items,
	}
	if len(items) != 1 {
		res.IsSequence = true
	}
	return res
}

func (ss Strings) MarshalYAML() (interface{}, error) {
	if ss.IsSequence || len(ss.Items) != 1 {
		return ss.Items, nil
	}
	return ss.Items[0], nil
}

func (ss *Strings) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if node.Decode(&s) == nil {
		ss.Items = []string{s}
		return nil
	}
	ss.IsSequence = true
	return node.Decode(&ss.Items)
}

func (ss Strings) Slice() []string {
	res := make([]string, len(ss.Items))
	for i, s := range ss.Items {
		res[i] = s
	}
	return res
}
