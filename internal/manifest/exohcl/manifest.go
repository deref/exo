package exohcl

import (
	"github.com/hashicorp/hcl/v2"
)

type Manifest struct {
	// Analysis inputs.
	Filename string
	File     *hcl.File

	// Analysis outputs.
	Content       *hcl.BodyContent
	FormatVersion *FormatVersion
	Environment   hcl.Blocks
	Components    hcl.Blocks
}

func NewManifest(filename string, file *hcl.File) *Manifest {
	return &Manifest{
		Filename: filename,
		File:     file,
	}
}

func (m *Manifest) Analyze(ctx *AnalysisContext) {
	var diags hcl.Diagnostics
	m.Content, diags = m.File.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "exo", Required: true},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "environment"},
			{Type: "components"},
		},
	})
	ctx.AppendDiags(diags...)

	if m.Content != nil {
		m.FormatVersion = NewFormatVersion(m)
		m.FormatVersion.Analyze(ctx)
		m.Environment = m.Content.Blocks.OfType("environment")
		m.Components = m.Content.Blocks.OfType("components")
	}
}
