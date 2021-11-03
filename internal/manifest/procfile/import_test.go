package procfile

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	// TODO: Improve formatting behavior.
	testImport(t, "parsed_call", `
one: A=1 two three four
`, `
exo = "0.1"
components {
  process "one" {
    program     = "two"
    arguments   = ["three", "four"]
    environment = { A = "1", PORT = "5000" }
  }
}
`)
	testImport(t, "shell form", `
thing: if make think; then X=1 ./thing "$@"; fi
chain: true && thing
`, `
exo = "0.1"
components {
  process "thing" {
    program     = "/bin/sh"
    arguments   = ["-c", "if make think; then X=1 ./thing \"$@\"; fi"]
    environment = { PORT = "5000" }
  }
  process "chain" {
    program     = "/bin/sh"
    arguments   = ["-c", "true && thing"]
    environment = { PORT = "5100" }
  }
}
`)
}

func testImport(t *testing.T, name, procfile, hcl string) {
	t.Run(name, func(t *testing.T) {
		ctx := context.Background()
		importer := &Importer{}
		analysisContext := &exohcl.AnalysisContext{
			Context: ctx,
		}
		f := importer.Import(analysisContext, []byte(procfile))
		diags := analysisContext.Diagnostics
		if len(diags) > 0 {
			t.Error(diags)
			return
		}
		var buf bytes.Buffer
		_, err := hclgen.WriteTo(&buf, hclgen.FileFromStructure(f))
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, strings.TrimSpace(hcl), string(bytes.TrimSpace(buf.Bytes())))
	})
}
