package compose

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration time.Duration

func (d Duration) String() string {
	td := time.Duration(d)
	var buf strings.Builder

	step := func(unit string, resolution time.Duration) {
		truncated := td.Truncate(resolution)
		td -= truncated
		if truncated > 0 {
			scale := truncated / resolution
			_, _ = fmt.Fprintf(&buf, "%d%s", scale, unit)
		}
	}
	step("h", time.Hour)
	step("m", time.Minute)
	step("s", time.Second)
	step("ms", time.Millisecond)
	step("us", time.Microsecond)

	if buf.Len() == 0 {
		return "0s"
	}
	return buf.String()
}

func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	var err error
	*d, err = ParseDuration(s)
	return err
}

func ParseDuration(s string) (Duration, error) {
	td, err := time.ParseDuration(s)
	return Duration(td), err
}
