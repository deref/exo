package compose

import "gopkg.in/yaml.v3"

// A sequence of strings that may be marshalled as an individual string.
type Tuple struct {
	IsSequence bool
	Items      []string
}

func MakeTuple(items ...string) Tuple {
	res := Tuple{
		Items: items,
	}
	if len(items) != 1 {
		res.IsSequence = true
	}
	return res
}

func (ss Tuple) MarshalYAML() (interface{}, error) {
	if ss.IsSequence || len(ss.Items) != 1 {
		return ss.Items, nil
	}
	return ss.Items[0], nil
}

func (ss *Tuple) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if node.Decode(&s) == nil {
		ss.Items = []string{s}
		return nil
	}
	ss.IsSequence = true
	return node.Decode(&ss.Items)
}

func (ss Tuple) Slice() []string {
	res := make([]string, len(ss.Items))
	for i, s := range ss.Items {
		res[i] = s
	}
	return res
}
