package compose

import (
	"errors"
	"strconv"
	"strings"

	"github.com/deref/exo/internal/providers/docker/compose/template"
)

type Bytes struct {
	Value      int64
	Expression string
}

func (bs Bytes) MarshalYAML() (interface{}, error) {
	i, err := strconv.ParseInt(bs.Expression, 10, 64)
	if err != nil {
		return bs.Expression, nil
	}
	return i, nil
}

func (bs *Bytes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&bs.Expression)
}

func (bs *Bytes) Interpolate(env Environment) error {
	tmpl, err := template.New(bs.Expression)
	if err != nil {
		return err
	}
	str, err := template.Substitute(tmpl, env)
	if err != nil {
		return err
	}
	bs.Value, err = ParseBytes(str)
	return err
}

var byteUnits = []byteUnit{
	{"gb", 1024 * 1024 * 1024},
	{"mb", 1024 * 1024},
	{"kb", 1024},
	{"g", 1024 * 1024 * 1024},
	{"m", 1024 * 1024},
	{"k", 1024},
	{"b", 1},
	{"", 1},
}

type byteUnit struct {
	Suffix string
	Scalar int64
}

func ParseBytes(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	s = strings.ToLower(s)
	for _, unit := range byteUnits {
		if strings.HasSuffix(s, unit.Suffix) {
			digits := s[:len(s)-len(unit.Suffix)]
			n, err := strconv.ParseInt(digits, 10, 64)
			if err != nil {
				break
			}
			n *= unit.Scalar
			return n, nil
		}
	}
	return 0, errors.New("expected integer number of bytes with a b, k, m, or g units suffix")
}
