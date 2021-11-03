package procfile

import (
	"bytes"
	"sort"
	"strconv"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const BasePort = 5000
const PortStep = 100

type Importer struct{}

func (imp *Importer) Import(ctx *exohcl.AnalysisContext, bs []byte) *hcl.File {
	b := exohcl.NewBuilder(bs)

	procfile, diags := Parse(bytes.NewBuffer(bs))
	ctx.AppendDiags(diags...)
	if diags.HasErrors() {
		return b.Build()
	}

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
			var subject *hcl.Range
			ctx.AppendDiags(exohcl.NewRenameWarning(p.Name, name, subject))
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

			environmentKeys := make([]string, 0, len(environment))
			for k := range environment {
				environmentKeys = append(environmentKeys, k)
			}
			sort.Strings(environmentKeys)

			for _, k := range environmentKeys {
				envExpr.Items = append(envExpr.Items, hclsyntax.ObjectConsItem{
					KeyExpr:   hclgen.NewObjStringKey(k, p.Range),
					ValueExpr: hclgen.NewStringLiteral(environment[k], p.Range),
				})
			}
			attrs = append(attrs, &hclsyntax.Attribute{
				Name:     "environment",
				Expr:     envExpr,
				SrcRange: p.Range,
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
	return b.Build()
}
