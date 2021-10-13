package hclgen

import (
	"bytes"
	"testing"

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
			}
		}`,
	}
	for _, src := range tests {
		f, diags := hclsyntax.ParseConfig([]byte(src), "", hcl.InitialPos)
		if !assert.False(t, diags.HasErrors()) {
			continue
		}
		var buf bytes.Buffer
		_, err := WriteTo(&buf, &hclsyntax.File{
			Body:  f.Body.(*hclsyntax.Body),
			Bytes: f.Bytes,
		})
		if !assert.NoError(t, err) {
			continue
		}
		formatted := string(hclwrite.Format([]byte(src)))

		cleaned := buf.Bytes()
		cleaned = hclwrite.Format(cleaned)
		cleaned = bytes.TrimSpace(cleaned)
		cleaned = bytes.ReplaceAll(cleaned, []byte{'\t'}, []byte("  "))

		assert.Equal(t, formatted, string(cleaned))
	}
}
