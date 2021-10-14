package exohcl

import (
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
		Attributes: hclgen.NewAttributes(&hclsyntax.Attribute{
			Name: "exo",
			Expr: hclgen.NewStringLiteral(Latest.String(), hcl.Range{}),
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
