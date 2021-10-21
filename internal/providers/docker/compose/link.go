package compose

import (
	"errors"
	"strings"

	"github.com/deref/exo/internal/providers/docker/compose/template"
	"gopkg.in/yaml.v3"
)

type Links []Link

func (ls Links) Values() []string {
	res := make([]string, len(ls))
	for i, link := range ls {
		res[i] = link.String.Value
	}
	return res
}

type Link struct {
	String
	Service string
	Alias   string
}

func (l Link) MarshalYAML() (interface{}, error) {
	return l.String, nil
}

func (l *Link) UnmarshalYAML(node *yaml.Node) error {
	if err := node.Decode(&l.String); err != nil {
		return err
	}
	_ = l.Interpolate(template.ErrEnvironment)
	return nil
}

func (l *Link) Interpolate(env Environment) error {
	if err := l.String.Interpolate(env); err != nil {
		return err
	}
	parts := strings.Split(l.String.Value, ":")
	switch len(parts) {
	case 1:
		l.Service = parts[0]
		l.Alias = parts[0]
	case 2:
		l.Service = parts[0]
		l.Alias = parts[1]
	default:
		return errors.New("expected SERVICE or SERVICE:ALIAS")
	}
	return nil
}
