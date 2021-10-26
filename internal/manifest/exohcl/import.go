package exohcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Importer struct {
	Filename string
}

func (imp *Importer) Import(ctx *AnalysisContext, bs []byte) *hcl.File {
	file, diags := hclsyntax.ParseConfig(bs, imp.Filename, hcl.InitialPos)
	ctx.AppendDiags(diags...)
	return file
}
