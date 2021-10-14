package testutil

import (
	"bytes"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func CleanHCL(bs []byte) string {
	bs = hclwrite.Format(bs)
	bs = bytes.TrimSpace(bs)
	bs = bytes.ReplaceAll(bs, []byte{'\t'}, []byte("  "))
	return string(bs)
}
