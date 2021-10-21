package template

import (
	"errors"
	"fmt"
)

type Environment interface {
	Lookup(*Variable) (string, error)
}

type MapEnvironment map[string]string

func (env MapEnvironment) Lookup(v *Variable) (substitution string, err error) {
	value, found := env[v.Name]
	switch v.Separator {
	case "":
		// No-op.
	case ":-":
		if value == "" {
			return v.DefaultOrError, nil
		}
	case "-":
		if !found {
			return v.DefaultOrError, nil
		}
	case ":?":
		if value == "" {
			return "", errors.New(v.DefaultOrError)
		}
	case "?":
		if !found {
			return "", errors.New(v.DefaultOrError)
		}
	default:
		return "", fmt.Errorf("invalid variable syntax: %q", v.Name+v.Separator)
	}
	return value, nil
}

type errEnvironment struct{}

var ErrEnvironment Environment = errEnvironment{}

func (_ errEnvironment) Lookup(v *Variable) (string, error) {
	return "", fmt.Errorf("undefined: %q", v.Name)
}
