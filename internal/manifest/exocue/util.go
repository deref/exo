package exocue

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
)

func dumpValue(v cue.Value) {
	bs, err := format.Node(v.Syntax())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bs)
}
