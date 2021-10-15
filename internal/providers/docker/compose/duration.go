package compose

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration struct {
	String
	Duration time.Duration
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	if err := node.Decode(&d.String); err != nil {
		return err
	}

	n, err := strconv.Atoi(d.Value)
	if err == nil {
		d.Duration = time.Duration(n) * time.Microsecond
		return nil
	}

	d.Duration, err = time.ParseDuration(d.Value)
	return err
}

func (d Duration) MarshalYAML() (interface{}, error) {
	if d.String.Expression != "" {
		return d.String, nil
	}

	td := d.Duration
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
		return "0s", nil
	}
	return buf.String(), nil
}
