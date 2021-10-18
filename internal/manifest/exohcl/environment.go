package exohcl

import (
	"fmt"

	"github.com/deref/exo/internal/environment"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Environment struct {
	m      *Manifest
	blocks hcl.Blocks
	vars   map[string]string
}

func newEnvironment(m *Manifest, blocks hcl.Blocks) *Environment {
	env := &Environment{
		m:      m,
		blocks: blocks,
		vars:   make(map[string]string),
	}

	if len(blocks) > 1 {
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Expected at most one environment block",
			Detail:   fmt.Sprintf("Only one environment block may appear in a manifest, but found %d", len(blocks)),
			Subject:  blocks[1].DefRange.Ptr(),
		})
	}

	for _, block := range blocks {
		if len(block.Labels) > 0 {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected label on environment block",
				Detail:   fmt.Sprintf("A environment block expects no labels, but has %d", len(block.Labels)),
				Subject:  &block.LabelRanges[0],
			})
		}

		attrs, diags := block.Body.JustAttributes()
		m.appendDiags(diags...)

		if !diags.HasErrors() {
			for _, attr := range attrs {
				// TODO: Validate attribute name.
				v, diags := attr.Expr.Value(env.m.evalCtx)
				m.appendDiags(diags...)
				if diags.HasErrors() {
					continue
				}
				if v.Type() != cty.String {
					m.appendDiags(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "expected environment variable to be a string",
						Detail:   fmt.Sprintf("environment variable evaluated to %s, but must be a string", v.Type().FriendlyName()),
						Subject:  attr.Expr.Range().Ptr(),
						Context:  &attr.Range,
					})
				}
				env.vars[attr.Name] = v.AsString()
			}
		}
	}

	return env
}

func (env *Environment) EnvironmentSource() string {
	return "manifest"
}

func (env *Environment) ExtendEnvironment(b environment.Builder) error {
	for k, v := range env.vars {
		b.AppendVariable(env, k, v)
	}
	return nil
}
