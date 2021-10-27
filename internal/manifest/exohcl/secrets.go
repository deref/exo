package exohcl

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Secrets struct {
	// Analysis inputs.
	Block *hclsyntax.Block

	// Analysis outputs.
	Source string
}

func NewSecrets(block *hclsyntax.Block) *Secrets {
	return &Secrets{
		Block: block,
	}
}

func (s *Secrets) Analyze(ctx *AnalysisContext) {
	if len(s.Block.Labels) > 0 {
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected label on secrets block",
			Detail:   fmt.Sprintf("A secrets block expects no labels, but has %d", len(s.Block.Labels)),
			Subject:  &s.Block.LabelRanges[0],
		})
	}

	content, diags := s.Block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "source", Required: true},
		},
	})
	ctx.AppendDiags(diags...)

	sourceAttr := content.Attributes["source"]
	s.Source, _ = AnalyzeString(ctx, sourceAttr.Expr)
}

type AppendSecrets struct {
	context.Context
	RewriteBase

	Source string
}

func (args AppendSecrets) RewriteEnvironment(re Rewrite, env *Environment) *hclgen.Block {
	res := args.RewriteBase.RewriteEnvironment(re, env)
	if res == nil {
		res = &hclgen.Block{
			Type: "environment",
			Body: &hclgen.Body{},
		}
	}
	res.Body.Blocks = append(res.Body.Blocks, &hclgen.Block{
		Type: "secrets",
		Body: &hclgen.Body{
			Attributes: []*hclsyntax.Attribute{
				{
					Name: "source",
					Expr: hclgen.NewStringLiteral(args.Source, hcl.Range{}),
				},
			},
		},
	})
	return res
}

type RemoveSecrets struct {
	context.Context
	RewriteBase

	Source string
}

func (args RemoveSecrets) RewriteSecrets(re Rewrite, s *Secrets) *hclgen.Block {
	if s.Source == args.Source {
		return nil
	}
	return args.RewriteBase.RewriteSecrets(re, s)
}
