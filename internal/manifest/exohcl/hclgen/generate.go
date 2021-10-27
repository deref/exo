package hclgen

import (
	"bytes"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Writes an HCL file AST to w. Note that syntax trivia (comments and
// whitespace) are not preserved, and so this function is only appropriate for
// conversion/generation use cases.
func WriteTo(w io.Writer, f *File) (int64, error) {
	out := hclwrite.NewEmptyFile()
	genFileTo(out, f)
	return out.WriteTo(w)
}

func FormatFile(f *File) []byte {
	var buf bytes.Buffer
	_, err := WriteTo(&buf, f)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func FormatBlock(block *hcl.Block) []byte {
	f := hclwrite.NewEmptyFile()
	out := f.Body().AppendNewBlock(block.Type, block.Labels)
	genBodyTo(out.Body(), block.Body)
	return f.Bytes()
}

func FormatExpression(x hclsyntax.Expression) []byte {
	f := hclwrite.NewEmptyFile()
	f.Body().AppendUnstructuredTokens(TokensForExpression(x))
	return f.Bytes()
}

func genFileTo(out *hclwrite.File, in *File) {
	genBodyTo(out.Body(), in.Body)
}

func genBodyTo(out *hclwrite.Body, in hcl.Body) {
	body := BodyFromStructure(in)
	for _, attr := range body.Attributes {
		out.SetAttributeRaw(attr.Name, TokensForExpression(attr.Expr))
	}
	for _, block := range body.Blocks {
		out.AppendBlock(genBlock(block))
	}
}

func genBlock(in *Block) *hclwrite.Block {
	out := hclwrite.NewBlock(in.Type, in.Labels)
	genBodyTo(out.Body(), in.Body)
	return out
}
