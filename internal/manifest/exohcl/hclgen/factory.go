package hclgen

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

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

func NewIdentifier(name string, rng hcl.Range) *hclsyntax.ScopeTraversalExpr {
	return &hclsyntax.ScopeTraversalExpr{
		Traversal: []hcl.Traverser{
			hcl.TraverseRoot{
				Name:     name,
				SrcRange: rng,
			},
		},
		SrcRange: rng,
	}
}

func NewObjStringKey(name string, rng hcl.Range) hclsyntax.Expression {
	var expr hclsyntax.Expression
	if hclsyntax.ValidIdentifier(name) {
		expr = NewIdentifier(name, rng)
	} else {
		expr = NewStringLiteral(name, rng)
	}
	return &hclsyntax.ObjectConsKeyExpr{
		Wrapped: expr,
	}
}
