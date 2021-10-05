package compose

import (
	"strconv"
	"strings"
)

type BoolExpression string

func (x *BoolExpression) UnmarshalYAML(bs []byte) error {
	*x = BoolExpression(bs)
	return nil
}

func (x BoolExpression) MarshalYAML() (interface{}, error) {
	s := string(x)
	switch strings.ToLower(s) {
	case "y", "yes", "true", "on":
		return true, nil
	case "n", "no", "false", "off":
		return false, nil
	default:
		return s, nil
	}
}

type Int64Expression string

func (x *Int64Expression) UnmarshalYAML(bs []byte) error {
	*x = Int64Expression(bs)
	return nil
}

func (x Int64Expression) MarshalYAML() (interface{}, error) {
	s := string(x)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return s, nil
	}
	return i, nil
}
