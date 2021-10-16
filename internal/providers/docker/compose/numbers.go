package compose

import (
	"strconv"

	"gopkg.in/yaml.v3"
)

type Int struct {
	String
	Value int64
}

func MakeInt(i int64) Int {
	s := strconv.FormatInt(i, 10)
	return Int{
		String: String{
			Tag:        "!!int",
			Expression: s,
			Value:      s,
		},
		Value: int64(i),
	}
}

func NewInt(i int64) *Int {
	ii := MakeInt(i)
	return &ii
}

func (i Int) Int() int {
	return int(i.Value)
}

func (i Int) Uint16() uint16 {
	return uint16(i.Value)
}

func (i *Int) Int64Ptr() *int64 {
	if i == nil {
		return nil
	}
	ii := i.Value
	return &ii
}

func (i *Int) UnmarshalYAML(node *yaml.Node) error {
	if err := i.String.UnmarshalYAML(node); err != nil {
		return err
	}
	_ = i.Interpolate(nil)
	return nil
}

func (i *Int) Interpolate(env Environment) error {
	if err := i.String.Interpolate(env); err != nil {
		return err
	}
	var err error
	i.Value, err = strconv.ParseInt(i.String.Value, 10, 64)
	return err
}
