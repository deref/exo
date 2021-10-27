package hclgen

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func FileEquiv(a, b *hcl.File) bool {
	return BodyEquiv(a.Body, b.Body)
}

func BodyEquiv(a, b hcl.Body) bool {
	synA := BodyFromStructure(a).SyntaxBody()
	synB := BodyFromStructure(b).SyntaxBody()
	return AttributesEquiv(synA.Attributes, synB.Attributes) && BlocksEquiv(synA.Blocks, synB.Blocks)
}

func AttributesEquiv(a, b hclsyntax.Attributes) bool {
	if len(a) != len(b) {
		return false
	}
	for k, attrA := range a {
		attrB, ok := b[k]
		if !ok {
			return false
		}
		if !AttributeEquiv(attrA, attrB) {
			return false
		}
	}
	return true
}

func AttributeEquiv(a, b *hclsyntax.Attribute) bool {
	return a.Name == b.Name && ExpressionEquiv(a.Expr, b.Expr)
}

func BlocksEquiv(a, b hclsyntax.Blocks) bool {
	n := len(a)
	if len(b) != n {
		return false
	}
	for i := 0; i < n; i++ {
		if !BlockEquiv(a[i], b[i]) {
			return false
		}
	}
	return true
}

func BlockEquiv(a, b *hclsyntax.Block) bool {
	return a.Type == b.Type && LabelsEquiv(a.Labels, b.Labels) && BodyEquiv(a.Body, b.Body)
}

func LabelsEquiv(a, b []string) bool {
	n := len(a)
	if len(b) != n {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ExpressionEquiv(a, b hcl.Expression) bool {
	ctx := &hcl.EvalContext{}
	va, diagsA := a.Value(ctx)
	vb, diagsB := b.Value(ctx)
	if diagsA.HasErrors() || diagsB.HasErrors() {
		panic("error evaluating expressions")
	}
	return va.Equals(vb).True()
}
