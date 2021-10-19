// This file contains utility functions for analyzing an Exo HCL AST.

package exohcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

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
