package compose

import "gopkg.in/yaml.v3"

// A sequence of strings that may be marshalled as an individual string.
type Tuple struct {
	IsSequence bool
	Items      []String
}

func MakeTuple(items ...String) Tuple {
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
	var s String
	if node.Decode(&s) == nil {
		ss.Items = []String{s}
		return nil
	}
	ss.IsSequence = true
	return node.Decode(&ss.Items)
}

func (tup *Tuple) Interpolate(env Environment) error {
	return interpolateSlice(tup.Items, env)
}

func (ss Tuple) Values() []string {
	res := make([]string, len(ss.Items))
	for i, s := range ss.Items {
		res[i] = s.Value
	}
	return res
}
