package hclgen_test

import (
	"bytes"
	"testing"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := []string{
		``,
		`x = 1`,
		`block {
		}`,
		`parent {
			a = 1
			b = 2
			child "label" {
				c = "d"
			}
			child "p" "q" {
				e = f(x)
				obj = { bare = "BARE", "quoted" = "QUOTED" }
			}
		}`,
	}
	for _, src := range tests {
		f, diags := hclsyntax.ParseConfig([]byte(src), "", hcl.InitialPos)
		if !assert.False(t, diags.HasErrors()) {
			continue
		}
		bs := hclgen.FormatFile(&hcl.File{
			Body:  f.Body.(*hclsyntax.Body),
			Bytes: f.Bytes,
		})
		formatted := string(hclwrite.Format([]byte(src)))

		bs = hclwrite.Format(bs)
		bs = bytes.TrimSpace(bs)
		bs = bytes.ReplaceAll(bs, []byte{'\t'}, []byte("  "))

		assert.Equal(t, formatted, string(bs))
	}
}
