package exohcl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	ctyyaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

type Manifest struct {
	filename    string
	f           *hcl.File
	evalCtx     *hcl.EvalContext
	diags       hcl.Diagnostics
	content     *hcl.BodyContent
	environment *Environment
	components  *ComponentSet
}

func Parse(filename string, bs []byte) *Manifest {
	file, diags := hclsyntax.ParseConfig(bs, filename, hcl.InitialPos)
	return NewManifest(filename, file, diags)
}

func NewManifest(filename string, f *hcl.File, diags hcl.Diagnostics) *Manifest {
	m := &Manifest{
		filename: filename,
		f:        f,
		diags:    diags,
		evalCtx: &hcl.EvalContext{
			Functions: map[string]function.Function{
				"jsonencode": stdlib.JSONEncodeFunc,
				"yamlencode": ctyyaml.YAMLEncodeFunc,
			},
		},
	}

	m.content, diags = m.f.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "exo", Required: true},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "environment"},
			{Type: "components"},
		},
	})
	m.appendDiags(diags...)

	var environmentBlocks, componentBlocks hcl.Blocks
	if m.content != nil {
		environmentBlocks = m.content.Blocks.OfType("environment")
		componentBlocks = m.content.Blocks.OfType("components")
	}
	m.environment = newEnvironment(m, environmentBlocks)
	m.components = newComponentSet(m, componentBlocks)

	return m
}

func (m *Manifest) appendDiags(diags ...*hcl.Diagnostic) {
	m.diags = append(m.diags, diags...)
}

func (m *Manifest) Diagnostics() hcl.Diagnostics {
	return m.diags
}

func (m *Manifest) File() *hcl.File {
	return m.f
}

func (m *Manifest) FormatVersion() (version FormatVersion) {
	attr := m.content.Attributes["exo"]
	s, diag := parseLiteralString(attr.Expr)
	if diag != nil {
		m.appendDiags(diag)
		return
	}
	parts := strings.Split(s, ".")
	ok := len(parts) == 2
	ints := make([]int, len(parts))
	for index, part := range parts {
		i, err := strconv.Atoi(part)
		ok = ok && err == nil && i >= 0
		ints[index] = i
	}
	if !ok {
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid exo format version constraint",
			Detail:   fmt.Sprintf(`Exo format constraint should be specified as 'major.minor'. The latest version is %q.`, Latest),
			Subject:  attr.Expr.Range().Ptr(),
			Context:  attr.Range.Ptr(),
		})
		return
	}
	version.Major = ints[0]
	version.Minor = ints[1]
	switch version.Major {
	case 1:
		if Latest.Minor < version.Minor {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported exo format minor version",
				Detail:   fmt.Sprintf(`Unsupported exo format minor version. The maximum supported "1.x" version is %q.`, Latest),
				Subject:  attr.Expr.Range().Ptr(),
				Context:  attr.Range.Ptr(),
			})
		}
	default:
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported exo format major version",
			Detail:   fmt.Sprintf(`Unsupported exo format major version. The latest version is %q.`, Latest),
			Subject:  attr.Expr.Range().Ptr(),
			Context:  attr.Range.Ptr(),
		})
	}
	return
}

func (m *Manifest) evalString(x hcl.Expression) string {
	v, diags := x.Value(m.evalCtx)
	m.appendDiags(diags...)
	if v.Type() != cty.String {
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected string",
			Detail:   fmt.Sprintf("Expected string, but got %s", v.Type().FriendlyName()),
			Subject:  x.Range().Ptr(),
		})
		return ""
	}
	return v.AsString()
}

func (m *Manifest) Environment() *Environment {
	return m.environment
}

func (m *Manifest) Components() *ComponentSet {
	return m.components
}
