package compose

import (
	"errors"
	"strconv"
	"strings"
)

type Bytes int64

func (bs Bytes) MarshalYAML() (interface{}, error) {
	return int64(bs), nil
}

func (bs *Bytes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var n int64
	err := unmarshal(&n)
	if err == nil {
		*bs = Bytes(n)
		return nil
	}
	var s string
	if err = unmarshal(&s); err != nil {
		return err
	}
	n, err = ParseBytes(s)
	*bs = Bytes(n)
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
