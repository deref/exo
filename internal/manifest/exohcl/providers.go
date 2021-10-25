package exohcl

import (
	"fmt"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type ProviderSet struct {
	m          *Manifest
	containers hcl.Blocks
	providers  []*Provider
}

func newProviderSet(m *Manifest, containers hcl.Blocks) *ProviderSet {
	ps := &ProviderSet{
		m:          m,
		containers: containers,
	}

	if len(containers) > 1 {
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Expected at most one providers block",
			Detail:   fmt.Sprintf("Only one providers block may appear in a manifest, but found %d", len(containers)),
			Subject:  containers[1].DefRange.Ptr(),
		})
	}

	for _, container := range containers {
		if len(container.Labels) > 0 {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected label on providers block",
				Detail:   fmt.Sprintf("A providers block expects no labels, but has %d", len(container.Labels)),
				Subject:  &container.LabelRanges[0],
			})
		}
		body, ok := container.Body.(*hclsyntax.Body)
		if !ok {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Malformed providers block",
				Detail:   fmt.Sprintf("Expected providers block to be an *hclsyntax.Body, but got %T", container.Body),
				Subject:  &container.DefRange,
			})
			continue
		}
		if len(body.Attributes) > 0 {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected attributes in providers block",
				Detail:   fmt.Sprintf("A providers block expects no attributes, but has %d", len(body.Attributes)),
				Subject:  body.Attributes.Range().Ptr(),
			})
		}
		ps.providers = make([]*Provider, len(body.Blocks))
		for i, block := range body.Blocks {
			ps.providers[i] = newProvider(m, block)
		}
	}

	return ps
}

func (ps *ProviderSet) Len() int {
	return len(ps.providers)
}

func (ps *ProviderSet) Index(i int) *Provider {
	return ps.providers[i]
}

type Provider struct {
	m         *Manifest
	source    *hclsyntax.Block
	expansion *hclsyntax.Block
	typ       string
	name      string
	options   map[string]string
}

func newProvider(m *Manifest, block *hclsyntax.Block) *Provider {
	p := &Provider{
		m:       m,
		source:  block,
		options: make(map[string]string),
	}

	if block.Type != "provider" {
		var diags hcl.Diagnostics
		block, diags = expandProvider(block)
		m.appendDiags(diags...)
	}
	p.expansion = block

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "type", Required: true},
			{Name: "options"},
		},
	})
	m.appendDiags(diags...)

	typeAttr := content.Attributes["type"]
	if typeAttr != nil {
		var diag *hcl.Diagnostic
		p.typ, diag = parseLiteralString(typeAttr.Expr)
		if diag != nil {
			m.appendDiags(diag)
		}
	}
	if p.typ == "" {
		p.typ = "invalid"
	}

	optionsAttr := content.Attributes["options"]
	if optionsAttr != nil {
		options, diags := optionsAttr.Expr.Value(m.evalCtx)
		m.appendDiags(diags...)
		if len(diags) == 0 {
			if !options.Type().IsMapType() {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Expected provider options to be a map.",
					Subject:  &optionsAttr.Range,
				})
			}
		}
		if len(diags) == 0 {
			valueMap := options.AsValueMap()
			for k, v := range valueMap {
				if v.Type() != cty.String {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  fmt.Sprintf("Expected provider option %q to have a string value.", k),
						Subject:  &optionsAttr.Range,
					})
					continue
				}
				p.options[k] = v.AsString()
			}
		}
	}

	switch len(block.Labels) {
	case 0:
		p.name = p.typ
	case 1:
		p.name = block.Labels[0]
	default:
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected label.",
			Detail:   `A provider block may have at most one label, which is overrides the name of the provider.`,
			Subject:  block.LabelRanges[1].Ptr(),
		})
	}

	return p
}

func (p *Provider) Options() map[string]string {
	res := make(map[string]string, len(p.options))
	for k, v := range p.options {
		res[k] = v
	}
	return res
}

func expandProvider(block *hclsyntax.Block) (*hclsyntax.Block, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	body := block.Body
	switch block.Type {
	case "esv":
		// OK.
	default:
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported provider type",
			Detail:   fmt.Sprintf(`The provider type %q is not recognized.`, block.Type),
			Subject:  block.DefRange().Ptr(),
		})
		return nil, diags
	}
	for _, subblock := range body.Blocks {
		switch subblock.Type {
		case "_":
			// TODO: copy content of "meta" blocks in to the expanded output.
		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected block",
				Detail:   fmt.Sprintf(`Unexpected %q block in %q provider.`, subblock.Type, block.Type),
				Subject:  body.Blocks[0].DefRange().Ptr(),
			})
		}
	}
	attrs := hclgen.AttributesFromSyntax(body.Attributes)
	optionItems := make([]hclsyntax.ObjectConsItem, 0, len(attrs))
	for _, attr := range attrs {
		optionItems = append(optionItems, hclsyntax.ObjectConsItem{
			KeyExpr:   hclgen.NewObjStringKey(attr.Name, attr.Range()),
			ValueExpr: attr.Expr,
		})
	}
	return &hclsyntax.Block{
		Type:   "provider",
		Labels: block.Labels,
		Body: &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{
				"type": &hclsyntax.Attribute{
					Name:        "type",
					Expr:        hclgen.NewStringLiteral(block.Type, block.TypeRange),
					SrcRange:    block.TypeRange,
					NameRange:   block.TypeRange,
					EqualsRange: block.TypeRange,
				},
				"options": &hclsyntax.Attribute{
					Name: "options",
					Expr: &hclsyntax.ObjectConsExpr{
						Items:     optionItems,
						SrcRange:  body.SrcRange,
						OpenRange: block.OpenBraceRange,
					},
					SrcRange:    body.SrcRange,
					NameRange:   block.TypeRange,
					EqualsRange: block.TypeRange,
				},
			},
		},
		TypeRange:       block.TypeRange,
		LabelRanges:     block.LabelRanges,
		OpenBraceRange:  block.OpenBraceRange,
		CloseBraceRange: block.CloseBraceRange,
	}, diags
}
