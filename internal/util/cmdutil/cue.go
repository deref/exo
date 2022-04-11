package cmdutil

import (
	"github.com/deref/exo/internal/util/cueutil"
)

func PrintCueStruct(v any) {
	Show(cueutil.MustValueToString(cueutil.EncodeValue(v)))
}
