package exohcl

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	ctx := context.Background()
	assertRewrite(t, Normalize{Context: ctx}, `
exo = "0.1"
`, `
exo = "0.1"
`)
}

func assertRewrite(t *testing.T, re Rewrite, before, after string) bool {
	ctx := context.Background()
	filename := "<file>"
	beforeFile, diags := hclsyntax.ParseConfig([]byte(before), filename, hcl.InitialPos)
	if len(diags) > 0 {
		t.Error(diags)
		return false
	}
	manifest := NewManifest(filename, beforeFile)
	if diags := Analyze(ctx, manifest); len(diags) > 0 {
		t.Error(diags)
		return false
	}
	afterFile := RewriteManifest(re, manifest)
	var buf bytes.Buffer
	hclgen.WriteTo(&buf, afterFile)
	expected := strings.TrimSpace(after)
	actual := string(bytes.TrimSpace(buf.Bytes()))
	return assert.Equal(t, expected, actual)
}
