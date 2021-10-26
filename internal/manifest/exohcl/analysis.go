// This file contains utility functions for analyzing an Exo HCL AST.
// Analysis is partially shallow/lazy. For deep/strict analysis, see Validate.

package exohcl

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	ctyyaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

type AnalysisContext struct {
	context.Context
	Diagnostics hcl.Diagnostics
}

func (ctx *AnalysisContext) AppendDiags(diags ...*hcl.Diagnostic) {
	ctx.Diagnostics = append(ctx.Diagnostics, diags...)
}

type Analyzer interface {
	Analyze(*AnalysisContext)
}

func Analyze(ctx context.Context, ana Analyzer) hcl.Diagnostics {
	analysisContext := &AnalysisContext{
		Context: ctx,
	}
	ana.Analyze(analysisContext)
	return analysisContext.Diagnostics
}

func parseLiteralString(x hcl.Expression) (string, *hcl.Diagnostic) {
	tmpl, ok := x.(*hclsyntax.TemplateExpr)
	if !(ok && tmpl.IsStringLiteral()) {
		return "", &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected literal string",
			Detail:   fmt.Sprintf("Expected literal string, got %T", x),
			Subject:  x.Range().Ptr(),
		}
	}
	lit := tmpl.Parts[0].(*hclsyntax.LiteralValueExpr)
	if lit.Val.Type() != cty.String {
		return "", &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected literal string",
			Detail:   fmt.Sprintf("Expected literal string, got %s", lit.Val.Type().FriendlyName()),
			Subject:  x.Range().Ptr(),
		}
	}
	return lit.Val.AsString(), nil
}

func NewRenameWarning(originalName, newName string, subject *hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  fmt.Sprintf("invalid name: %s, renamed to %q", originalName, newName),
		Subject:  subject,
	}
}

func NewUnsupportedFeatureWarning(featureName, explanation string, subject *hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  fmt.Sprintf("unsupported feature: %s", featureName),
		Detail:   fmt.Sprintf("The %s feature is unsupported. %s", featureName, explanation),
		Subject:  subject,
	}
}

var evalCtx = &hcl.EvalContext{
	Functions: map[string]function.Function{
		"jsonencode": stdlib.JSONEncodeFunc,
		"yamlencode": ctyyaml.YAMLEncodeFunc,
	},
}

func AnalyzeString(ctx *AnalysisContext, x hcl.Expression) (s string, ok bool) {
	v, diags := x.Value(evalCtx)
	ctx.AppendDiags(diags...)
	if v.Type() != cty.String {
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected string",
			Detail:   fmt.Sprintf("Expected string, but got %s", v.Type().FriendlyName()),
			Subject:  x.Range().Ptr(),
		})
		return "", false
	}
	return v.AsString(), true
}
