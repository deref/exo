package cmdutil

import (
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"github.com/deref/exo/internal/util/cueutil"
)

func PrintCueStruct(v any) {
	cc := cuecontext.New()
	value := cc.Encode(v)
	bs, err := format.Node(cueutil.StructToFile(value))
	if err != nil {
		panic(err)
	}
	Show(string(bs))
}
