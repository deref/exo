// This file contains AST nodes that preserve attribute order.

package hclgen

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Block struct {
	Type   string
	Labels []string
	Body   *Body

	TypeRange       hcl.Range
	LabelRanges     []hcl.Range
	OpenBraceRange  hcl.Range
	CloseBraceRange hcl.Range
}

func BlocksFromSyntax(blocks []*hclsyntax.Block) []*Block {
	res := make([]*Block, len(blocks))
	for i, block := range blocks {
		res[i] = BlockFromSyntax(block)
	}
	return res
}

func BlockFromSyntax(in *hclsyntax.Block) *Block {
	return &Block{
		Type:   in.Type,
		Labels: in.Labels,
		Body:   BodyFromSyntax(in.Body),

		TypeRange:       in.TypeRange,
		LabelRanges:     in.LabelRanges,
		OpenBraceRange:  in.OpenBraceRange,
		CloseBraceRange: in.CloseBraceRange,
	}
}

func (b *Block) SyntaxBlock() *hclsyntax.Block {
	return &hclsyntax.Block{
		Type:   b.Type,
		Labels: b.Labels,
		Body:   b.Body.SyntaxBody(),

		TypeRange:       b.TypeRange,
		LabelRanges:     b.LabelRanges,
		OpenBraceRange:  b.OpenBraceRange,
		CloseBraceRange: b.CloseBraceRange,
	}
}

type Body struct {
	Attributes []*hclsyntax.Attribute
	Blocks     []*Block

	SrcRange hcl.Range
	EndRange hcl.Range
}

func BodyFromStructure(body hcl.Body) *Body {
	switch body := body.(type) {
	case *hclsyntax.Body:
		return BodyFromSyntax(body)
	case *Body:
		return body
	default:
		panic(fmt.Errorf("unsupported body type: %T", body))
	}
}

func BodyFromSyntax(body *hclsyntax.Body) *Body {
	return &Body{
		Attributes: AttributesFromSyntax(body.Attributes),
		Blocks:     BlocksFromSyntax(body.Blocks),

		SrcRange: body.SrcRange,
		EndRange: body.EndRange,
	}
}

func (b *Body) SyntaxAttributes() hclsyntax.Attributes {
	res := make(hclsyntax.Attributes)
	for _, attr := range b.Attributes {
		res[attr.Name] = attr
	}
	return res
}

func (b *Body) SyntaxBlocks() []*hclsyntax.Block {
	res := make([]*hclsyntax.Block, len(b.Blocks))
	for i, block := range b.Blocks {
		res[i] = block.SyntaxBlock()
	}
	return res
}

func (b *Body) SyntaxBody() *hclsyntax.Body {
	return &hclsyntax.Body{
		Attributes: b.SyntaxAttributes(),
		Blocks:     b.SyntaxBlocks(),
		SrcRange:   b.SrcRange,
		EndRange:   b.EndRange,
	}
}

func (b *Body) Content(schema *hcl.BodySchema) (*hcl.BodyContent, hcl.Diagnostics) {
	return b.SyntaxBody().Content(schema)
}

func (b *Body) PartialContent(schema *hcl.BodySchema) (*hcl.BodyContent, hcl.Body, hcl.Diagnostics) {
	return b.SyntaxBody().PartialContent(schema)
}

func (b *Body) JustAttributes() (hcl.Attributes, hcl.Diagnostics) {
	return b.SyntaxBody().JustAttributes()
}

func (b *Body) MissingItemRange() hcl.Range {
	return b.SyntaxBody().MissingItemRange()
}

type Attributes []*Attribute

type Attribute = hclsyntax.Attribute

func AttributesFromSyntax(attributes hclsyntax.Attributes) Attributes {
	res := make(Attributes, 0, len(attributes))
	for _, attribute := range attributes {
		res = append(res, attribute)
	}
	sort.Sort(res)
	return res
}

func (a Attributes) Len() int      { return len(a) }
func (a Attributes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Try to preserve the source code order, but fallback to alphabetical.
func (a Attributes) Less(i, j int) bool {
	lhs := a[i]
	rhs := a[j]
	if lhs.Range().Start.Byte < rhs.Range().Start.Byte {
		return true
	}
	return lhs.Name < rhs.Name
}
