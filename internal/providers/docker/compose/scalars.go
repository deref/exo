package compose

import (
	"fmt"
	"strconv"
	"strings"

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

type Bool struct {
	Value      bool
	Expression string
}

func (b *Bool) Interpolate(env Environment) error {
	tmpl, err := template.New(b.Expression)
	if err != nil {
		return err
	}
	s, err := template.Substitute(tmpl, env)
	if err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "y", "yes", "true", "on":
		b.Value = true
	case "n", "no", "false", "off", "":
		b.Value = false
	default:
		return fmt.Errorf("invalid boolean: %q", s)
	}
	return nil
}

type Int struct {
	Value      int
	Expression string
}

func (i *Int) Interpolate(env Environment) error {
	tmpl, err := template.New(i.Expression)
	if err != nil {
		return err
	}
	s, err := template.Substitute(tmpl, env)
	if err != nil {
		return err
	}
	i.Value, err = strconv.Atoi(s)
	return err
}
