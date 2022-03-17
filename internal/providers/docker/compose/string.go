package compose

import (
	"fmt"

	"github.com/deref/exo/internal/providers/docker/compose/template"
	"gopkg.in/yaml.v3"
)

type Strings []String

func (ss Strings) Values() []string {
	res := make([]string, len(ss))
	for i, s := range ss {
		res[i] = s.Value
	}
	return res
}

// String is a scalar leaf node in a compose document.  Non-string types
// generally embed a String to facilitate interpolation.  See the package
// documentation.
type String struct {
	Tag        string
	Style      yaml.Style
	Expression string
	Value      string
}

func MakeString(s string) String {
	return String{
		Tag:        "!!str",
		Expression: s,
		Value:      s,
	}
}

func (s String) WithValue(v string) String {
	s.Value = v
	return s
}

func (s *String) UnmarshalYAML(node *yaml.Node) error {
	s.Tag = node.Tag
	s.Style = node.Style
	err := node.Decode(&s.Expression)
	s.Value = s.Expression
	_ = s.Interpolate(ErrEnvironment)
	return err
}

func (s String) MarshalYAML() (any, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   s.Tag,
		Style: s.Style,
		Value: s.Expression,
	}, nil
}

func (s *String) Interpolate(env Environment) error {
	tmpl, err := template.Parse(s.Expression)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	s.Value, err = template.Substitute(tmpl, env)
	return err
}
