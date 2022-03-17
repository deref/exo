// This file implements HTTPie-like parsing of JSON objects from command line
// args.

package cmdutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func ArgsToJsonObject(args []string) (map[string]any, error) {
	m := make(map[string]any, len(args))
	for _, arg := range args {
		k, v, err := ArgToJsonEntry(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid json entry %q: %w", arg, err)
		}
		m[k] = v
	}
	return m, nil
}

func ArgToJsonEntry(arg string) (k string, v any, err error) {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		err = errors.New(`expected "=" or ":="`)
		return
	}
	lhs := parts[0]
	rhs := parts[1]
	if strings.HasSuffix(lhs, ":") {
		k = strings.TrimSuffix(lhs, ":")
		dec := json.NewDecoder(strings.NewReader(rhs))
		err = dec.Decode(&v)
	} else {
		k = lhs
		v = rhs
	}
	return
}
