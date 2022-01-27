package exocue

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
)

func FormatBytes(v cue.Value) ([]byte, error) {
	return format.Node(v.Syntax())
}

func FormatString(v cue.Value) (string, error) {
	bs, err := FormatBytes(v)
	return string(bs), err
}

func dumpValue(v cue.Value) {
	bs, err := FormatBytes(v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bs)
}
