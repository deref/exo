package compose

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Bool struct {
	String
	Value bool
}

func MakeBool(b bool) Bool {
	var s string
	if b {
		s = "true"
	} else {
		s = "false"
	}
	return Bool{
		String: String{
			Tag:        "!!bool",
			Expression: s,
			Value:      s,
		},
		Value: b,
	}
}

func NewBool(b bool) *Bool {
	bb := MakeBool(b)
	return &bb
}

func (b *Bool) Ptr() *bool {
	if b == nil {
		return nil
	}
	bb := b.Value
	return &bb
}

func (b *Bool) UnmarshalYAML(node *yaml.Node) error {
	if err := b.String.UnmarshalYAML(node); err != nil {
		return err
	}
	_ = b.Interpolate(ErrEnvironment)
	return nil
}

func (b *Bool) Interpolate(env Environment) error {
	if err := b.String.Interpolate(env); err != nil {
		return err
	}
	switch strings.ToLower(b.String.Value) {
	case "1", "y", "yes", "t", "true", "on":
		b.Value = true
	case "0", "n", "no", "f", "false", "off", "":
		b.Value = false
	default:
		return fmt.Errorf("malformed boolean: %q", b.String.Value)
	}
	return nil
}
