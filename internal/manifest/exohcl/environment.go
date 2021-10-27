package exohcl

import (
	"fmt"

	"github.com/deref/exo/internal/environment"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type Environment struct {
	// Analysis inputs.
	Blocks hcl.Blocks

	// Analysis outputs.
	Attributes []*hclgen.Attribute
	Variables  map[string]string
	Secrets    []*Secrets // TODO: Should analysis be lazier here?
}

func NewEnvironment(m *Manifest) *Environment {
	return &Environment{
		Blocks: m.Environment,
	}
}

func (env *Environment) Analyze(ctx *AnalysisContext) {
	env.Variables = make(map[string]string)

	if len(env.Blocks) > 1 {
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Expected at most one environment block",
			Detail:   fmt.Sprintf("Only one environment block may appear in a manifest, but found %d", len(env.Blocks)),
			Subject:  env.Blocks[1].DefRange.Ptr(),
		})
	}

	for _, block := range env.Blocks {
		if len(block.Labels) > 0 {
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected label on environment block",
				Detail:   fmt.Sprintf("A environment block expects no labels, but has %d", len(block.Labels)),
				Subject:  &block.LabelRanges[0],
			})
		}

		body := block.Body.(*hclsyntax.Body)

		for _, attr := range body.Attributes {
			// TODO: Validate attribute name.
			v, diags := attr.Expr.Value(evalCtx)
			ctx.AppendDiags(diags...)
			if diags.HasErrors() {
				continue
			}
			if v.Type() != cty.String {
				ctx.AppendDiags(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "expected environment variable to be a string",
					Detail:   fmt.Sprintf("environment variable evaluated to %s, but must be a string", v.Type().FriendlyName()),
					Subject:  attr.Expr.Range().Ptr(),
					Context:  &attr.SrcRange,
				})
			}
			env.Attributes = append(env.Attributes, attr)
			env.Variables[attr.Name] = v.AsString()
		}

		for _, child := range body.Blocks {
			switch child.Type {
			case "secrets":
				secrets := NewSecrets(child)
				secrets.Analyze(ctx)
				env.Secrets = append(env.Secrets, secrets)
			default:
				ctx.AppendDiags(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unexpected block in environment",
					Detail:   fmt.Sprintf("Unexpected %q block in environment block", child.Type),
					Subject:  &child.TypeRange,
					Context:  child.DefRange().Ptr(),
				})
			}
		}
	}
}

func (env *Environment) EnvironmentSource() string {
	return "manifest"
}

func (env *Environment) ExtendEnvironment(b environment.Builder) error {
	for k, v := range env.Variables {
		b.AppendVariable(env, k, v)
	}
	return nil
}
