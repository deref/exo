package procfile

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const BasePort = 5000
const PortStep = 100

type Converter struct{}

func (c *Converter) Convert(bs []byte) (*hcl.File, hcl.Diagnostics) {
	procfile, diags := Parse(bytes.NewBuffer(bs))
	if diags.HasErrors() {
		return nil, diags
	}

	b := exohcl.NewBuilder(bs)

	port := BasePort
	for _, p := range procfile.Processes {
		// Assign default PORT, merge in specified environment.
		environment := map[string]string{
			"PORT": strconv.Itoa(port),
		}
		for name, value := range p.Environment {
			environment[name] = value
		}
		port += PortStep

		// Get component name.
		name := exohcl.MangleName(p.Name)
		if name != p.Name {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  fmt.Sprintf("invalid name: %s, renamed to %q", p.Name, name),
			})
		}

		// Build HCL attributes.
		args := make([]hclsyntax.Expression, len(p.Arguments))
		for i, arg := range p.Arguments {
			args[i] = hclgen.NewStringLiteral(arg, p.Range)
		}
		attrs := []*hclsyntax.Attribute{
			{
				Name:     "program",
				Expr:     hclgen.NewStringLiteral(p.Program, p.Range),
				SrcRange: p.Range,
			},
			{
				Name:     "arguments",
				Expr:     hclgen.NewTuple(args, p.Range),
				SrcRange: p.CommandRange,
			},
		}
		if len(environment) > 0 {
			envExpr := &hclsyntax.ObjectConsExpr{
				SrcRange: p.Range,
			}
			for k, v := range environment {
				envExpr.Items = append(envExpr.Items, hclsyntax.ObjectConsItem{
					KeyExpr:   hclgen.NewObjStringKey(k, p.Range),
					ValueExpr: hclgen.NewStringLiteral(v, p.Range),
				})
			}
			attrs = append(attrs, &hclsyntax.Attribute{
				Name: "environment",
				Expr: envExpr,
			})
		}

		b.AddComponentBlock(&hclgen.Block{
			Type:   "process",
			Labels: []string{name},
			Body: &hclgen.Body{
				Attributes: attrs,
			},
		})
	}
	return b.Build(), diags
}
