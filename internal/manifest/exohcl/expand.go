package exohcl

import (
	"context"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
)

type Expand struct {
	context.Context
	RewriteBase
}

func (_ Expand) RewriteComponent(re Rewrite, c *Component) *hclgen.Block {
	return hclgen.BlockFromSyntax(c.Expansion)
}
