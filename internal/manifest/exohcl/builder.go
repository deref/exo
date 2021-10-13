package exohcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type Builder struct {
	file     *hcl.File
	topLevel *hclsyntax.Body
}

// If provided, src is the original source code used to build a new Exo
// manifest file. The actual format of the source code doesn't matter, since
// diagnostics will be reported in terms of this source code when available,
// not the constructed manifest file.
func NewBuilder(src []byte) *Builder {
	topLevel := &hclsyntax.Body{
		Attributes: NewAttributes(&hclsyntax.Attribute{
			Name: "exo",
			Expr: NewStringLiteral(Latest.String(), hcl.Range{}),
		}),
	}
	return &Builder{
		file: &hcl.File{
			Body:  topLevel,
			Bytes: src,
		},
		topLevel: topLevel,
	}
}

func (b *Builder) Build() *hcl.File {
	return b.file
}

func (b *Builder) EnsureComponentBlock() *hclsyntax.Block {
	for _, block := range b.topLevel.Blocks {
		if block.Type == "components" {
			return block
		}
	}
	block := &hclsyntax.Block{
		Type: "components",
		Body: &hclsyntax.Body{},
	}
	b.topLevel.Blocks = append(b.topLevel.Blocks, block)
	return block
}

func (b *Builder) AddComponentBlock(block *hclsyntax.Block) {
	componentsBlock := b.EnsureComponentBlock()
	componentsBlock.Body.Blocks = append(componentsBlock.Body.Blocks, block)
}

func NewAttributes(attributes ...*hclsyntax.Attribute) hclsyntax.Attributes {
	m := make(hclsyntax.Attributes, len(attributes))
	for _, attr := range attributes {
		m[attr.Name] = attr
	}
	return m
}

func NewNullLiteral(rng hcl.Range) *hclsyntax.LiteralValueExpr {
	return &hclsyntax.LiteralValueExpr{
		Val:      cty.NilVal,
		SrcRange: rng,
	}
}

func NewStringLiteral(s string, rng hcl.Range) *hclsyntax.TemplateExpr {
	return &hclsyntax.TemplateExpr{
		Parts: []hclsyntax.Expression{
			&hclsyntax.LiteralValueExpr{
				Val:      cty.StringVal(s),
				SrcRange: rng,
			},
		},
		SrcRange: rng,
	}
}

func NewTuple(exprs []hclsyntax.Expression, rng hcl.Range) *hclsyntax.TupleConsExpr {
	return &hclsyntax.TupleConsExpr{
		Exprs:    exprs,
		SrcRange: rng,
	}
}
