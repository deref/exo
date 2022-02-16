package cmdutil

import (
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"github.com/deref/exo/internal/manifest/exocue"
)

func PrintCueStruct(v interface{}) {
	cc := cuecontext.New()
	value := cc.Encode(v)
	bs, err := format.Node(exocue.StructToFile(value))
	if err != nil {
		panic(err)
	}
	Show(string(bs))
}
