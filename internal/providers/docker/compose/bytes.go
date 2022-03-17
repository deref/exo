package compose

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Bytes struct {
	String
	Quantity int64
	Unit     ByteUnit
}

type ByteUnit struct {
	Suffix string
	Scalar int64
}

func (bs Bytes) Int64() int64 {
	return bs.Quantity * bs.Unit.Scalar
}

func (bs Bytes) Uint64() uint64 {
	return uint64(bs.Int64())
}

func (bs Bytes) MarshalYAML() (any, error) {
	if bs.Unit.Suffix == "" {
		return bs.Quantity, nil
	}
	return fmt.Sprintf("%d%s", bs.Quantity, bs.Unit.Suffix), nil
}

func (bs *Bytes) UnmarshalYAML(node *yaml.Node) error {
	if err := node.Decode(&bs.String); err != nil {
		return err
	}
	_ = bs.Interpolate(ErrEnvironment)
	return nil
}

func (bs *Bytes) Interpolate(env Environment) error {
	if err := bs.String.Interpolate(env); err != nil {
		return err
	}
	return bs.Parse(bs.String.Value)
}

func (bs *Bytes) Parse(s string) error {
	if s == "" {
		return nil
	}

	n, err := strconv.ParseInt(bs.String.Value, 10, 64)
	if err == nil {
		bs.String.Tag = "!!int"
		bs.Quantity = n
		return nil
	}

	s = strings.ToLower(s)
	for _, unit := range byteUnits {
		if strings.HasSuffix(s, unit.Suffix) {
			digits := s[:len(s)-len(unit.Suffix)]
			quantity, err := strconv.ParseInt(digits, 10, 64)
			if err != nil {
				break
			}
			bs.Quantity = quantity
			bs.Unit = unit
			return nil
		}
	}
	return errors.New("expected integer number of bytes with a b, k, m, or g units suffix")
}

var byteUnits = []ByteUnit{
	{"gb", 1024 * 1024 * 1024},
	{"mb", 1024 * 1024},
	{"kb", 1024},
	{"g", 1024 * 1024 * 1024},
	{"m", 1024 * 1024},
	{"k", 1024},
	{"b", 1},
	{"", 1},
}
