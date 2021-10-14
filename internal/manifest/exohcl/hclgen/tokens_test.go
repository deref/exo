package hclgen

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestFormatExpression(t *testing.T) {
	exprs := []string{
		`1`,
		`true`,
		`false`,

		`"x${1}y"`,
		`"abc"`,

		`x`,
		`x.y.z`,
		`"${x}"`,
		`x["y"]`,
		// TODO: `x[*]`,

		`lhs || rhs`,
		`lhs && rhs`,
		`lhs == rhs`,
		`lhs != rhs`,
		`lhs < rhs`,
		`lhs <= rhs`,
		`lhs > rhs`,
		`lhs >= rhs`,
		`lhs + rhs`,
		`lhs - rhs`,
		`lhs * rhs`,
		`lhs / rhs`,
		`lhs % rhs`,

		`(x)`,
		`1 + 2 * 3`,
		`(1 + 2) * 3`,

		`!true`,
		`-1`,
		`--1`,

		`f()`,
		`f(1, 2, 3)`,

		`x ? y : z`,

		`[]`,
		`[1]`,
		`[1, 2, 3]`,

		`{}`,
		`{ x = 1 }`,
		`{ x = 1, y = 2 }`,
		`{ "x" = 1 }`,
		`{ 1 = 2 }`,

		// TODO: Test pretty-printing multi-line collections, heredocs, etc.
	}
	for _, src := range exprs {
		bs := []byte(src)
		ast, err := hclsyntax.ParseExpression(bs, "", hcl.InitialPos)
		if !assert.False(t, err.HasErrors()) {
			continue
		}
		assert.Equal(t, string(hclwrite.Format(bs)), string(FormatExpression(ast)))
	}
}
