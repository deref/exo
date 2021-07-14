package inflect

import (
	"bytes"
	"strings"
)

func KebabToPascal(s string) string {
	return KebabToGo(true, s)
}

func KebabToCamel(s string) string {
	return KebabToGo(false, s)
}

func KebabToGo(public bool, s string) string {
	var b strings.Builder
	b.Grow(len(s))
	up := public
	bs := []byte(s)
	for i, c := range bs {
		if c == '-' {
			up = true
			continue
		}
		cs := bs[i : i+1]
		if up {
			cs = bytes.ToUpper(cs)
			up = false
		}
		_, _ = b.Write(cs)
	}
	return b.String()
}
