package compose

import (
	"github.com/deref/exo/internal/providers/docker/compose/template"
)

type String struct {
	Value      string
	Expression string
}

func (s *String) Interpolate(env Environment) error {
	tmpl, err := template.New(s.Expression)
	if err != nil {
		return err
	}
	s.Value, err = template.Substitute(tmpl, env)
	return err
}
