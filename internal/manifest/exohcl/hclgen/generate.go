package hclgen

import (
	"io"
	"sort"

	"github.com/deref/exo/internal/manifest/exohcl/exohclwrite"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Writes an HCL file AST to w. Note that syntax trivia (comments and
// whitespace) are not preserved, and so this function is only appropriate for
// conversion/generation use cases.
func WriteTo(w io.Writer, f *hclsyntax.File) (int64, error) {
	out := hclwrite.NewEmptyFile()
	genFileTo(out, f)
	return out.WriteTo(w)
}

func genFileTo(out *hclwrite.File, in *hclsyntax.File) {
	genBodyTo(out.Body(), in.Body)
}

func genBodyTo(out *hclwrite.Body, in *hclsyntax.Body) {
	attrs := make([]*hclsyntax.Attribute, 0, len(in.Attributes))
	for _, attr := range in.Attributes {
		attrs = append(attrs, attr)
	}
	sort.Sort(attributeSort(attrs))
	for _, attr := range attrs {
		out.SetAttributeRaw(attr.Name, exohclwrite.TokensForExpression(attr.Expr))
	}
	for _, block := range in.Blocks {
		out.AppendBlock(genBlock(block))
	}
}

type attributeSort []*hclsyntax.Attribute

func (a attributeSort) Len() int      { return len(a) }
func (a attributeSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Try to preserve the source code order, but fallback to alphabetical.
func (a attributeSort) Less(i, j int) bool {
	lhs := a[i]
	rhs := a[j]
	if lhs.Range().Start.Byte < rhs.Range().Start.Byte {
		return true
	}
	return lhs.Name < rhs.Name
}

func genBlock(in *hclsyntax.Block) *hclwrite.Block {
	out := hclwrite.NewBlock(in.Type, in.Labels)
	genBodyTo(out.Body(), in.Body)
	return out
}
