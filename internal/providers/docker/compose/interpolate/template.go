package interpolate

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
)

type Template interface {
	Substitute(w io.Writer, env Environment) error
}

func Substitute(template Template, env Environment) (string, error) {
	var buf bytes.Buffer
	err := template.Substitute(&buf, env)
	return buf.String(), err
}

func NewTemplate(s string) (Template, error) {
	var elements []Template
	matches := pattern.FindAllStringSubmatchIndex(s, -1)
	left := 0
	for _, match := range matches {
		matchLeft := match[0]
		matchRight := match[1]
		escaped := match[2] != -1
		if escaped {
			elements = append(elements, &Literal{
				Value: s[left:matchLeft] + "$",
			})
			left = matchRight
			continue
		}

		if left < matchLeft {
			elements = append(elements, &Literal{
				Value: s[left:matchLeft],
			})
		}
		left = matchRight

		nameLeft := match[4]
		nameRight := match[5]
		if nameLeft == -1 {
			nameLeft = match[6]
			nameRight = match[7]
		}
		// End-of-string.
		if nameLeft == -1 {
			break
		}
		v := &Variable{
			Name: s[nameLeft:nameRight],
		}

		defaultLeft := match[10]
		defaultRight := match[11]
		if defaultLeft != -1 {
			v.DefaultOrError = s[defaultLeft:defaultRight]

			sepLeft := match[8]
			sepRight := match[9]
			v.Separator = s[sepLeft:sepRight]
		}

		elements = append(elements, v)
	}
	return &Sequence{
		Elements: elements,
	}, nil
}

var pattern = regexp.MustCompile("(?i)\\$(?:(\\$)|([_a-z][_a-z0-9]*)|{([_a-z][_a-z0-9]*)(?:(:?[-?])([^}]*))?})|$")

type Sequence struct {
	Elements []Template
}

func (seq *Sequence) Substitute(w io.Writer, env Environment) error {
	for _, element := range seq.Elements {
		if err := element.Substitute(w, env); err != nil {
			return err
		}
	}
	return nil
}

type Literal struct {
	Value string
}

func (lit *Literal) Substitute(w io.Writer, env Environment) error {
	_, err := io.WriteString(w, lit.Value)
	return err
}

type Variable struct {
	Name           string
	Separator      string
	DefaultOrError string
}

func (v *Variable) Substitute(w io.Writer, env Environment) error {
	value, found := env.Lookup(v.Name)
	switch v.Separator {
	case "":
		// No-op.
	case ":-":
		if value == "" {
			value = v.DefaultOrError
		}
	case "-":
		if !found {
			value = v.DefaultOrError
		}
	case ":?":
		if value == "" {
			return errors.New(v.DefaultOrError)
		}
	case "?":
		if !found {
			return errors.New(v.DefaultOrError)
		}
	default:
		return fmt.Errorf("invalid variable syntax: %q", v.Name+v.Separator)
	}
	_, err := io.WriteString(w, value)
	return err
}
