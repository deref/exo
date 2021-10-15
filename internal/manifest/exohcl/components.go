package exohcl

import (
	"fmt"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type ComponentSet struct {
	m          *Manifest
	containers hcl.Blocks
	components []*Component
}

func newComponentSet(m *Manifest, containers hcl.Blocks) *ComponentSet {
	cs := &ComponentSet{
		m:          m,
		containers: containers,
	}

	if len(containers) > 1 {
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Expected at most one components block",
			Detail:   fmt.Sprintf("Only one components block may appear in a manifest, but found %d", len(containers)),
			Subject:  containers[1].DefRange.Ptr(),
		})
	}

	for _, container := range containers {
		if len(container.Labels) > 0 {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected label on components block",
				Detail:   fmt.Sprintf("A components block expects no labels, but has %d", len(container.Labels)),
				Subject:  &container.LabelRanges[0],
			})
		}
		body, ok := container.Body.(*hclsyntax.Body)
		if !ok {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Malformed components block",
				Detail:   fmt.Sprintf("Expected components block to be an *hclsyntax.Body, but got %T", container.Body),
				Subject:  &container.DefRange,
			})
			continue
		}
		if len(body.Attributes) > 0 {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected attributes in components block",
				Detail:   fmt.Sprintf("A components block expects no attributes, but has %d", len(body.Attributes)),
				Subject:  body.Attributes.Range().Ptr(),
			})
		}
		cs.components = make([]*Component, len(body.Blocks))
		for i, block := range body.Blocks {
			cs.components[i] = newComponent(m, block)
		}
	}

	return cs
}

func (cs *ComponentSet) Len() int {
	return len(cs.components)
}

func (cs *ComponentSet) Index(i int) *Component {
	return cs.components[i]
}

type Component struct {
	m         *Manifest
	source    *hclsyntax.Block
	expansion *hclsyntax.Block
	typ       string
	name      string
	spec      string
	dependsOn []string
}

func newComponent(m *Manifest, block *hclsyntax.Block) *Component {
	c := &Component{
		m:      m,
		source: block,
	}

	if block.Type != "component" {
		var diags hcl.Diagnostics
		block, diags = expandComponent(block)
		m.appendDiags(diags...)
	}
	c.expansion = block

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "type", Required: true},
			{Name: "spec"},
			{Name: "depends_on"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "spec"},
		},
	})
	m.appendDiags(diags...)

	switch len(block.Labels) {
	case 0:
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected component name",
			Detail:   `A component block must have exactly one label, which is the name of the component.`,
			Subject:  block.DefRange().Ptr(),
		})
	case 1:
		c.name = block.Labels[0]
	default:
		m.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected label.",
			Detail:   `A component block must have exactly one label, which is the name of the component.`,
			Subject:  block.LabelRanges[1].Ptr(),
		})
	}

	typeAttr := content.Attributes["type"]
	if typeAttr != nil {
		var diag *hcl.Diagnostic
		c.typ, diag = parseLiteralString(typeAttr.Expr)
		if diag != nil {
			m.appendDiags(diag)
		}
	}

	specAttr := content.Attributes["spec"]
	if specAttr == nil {
		specBlocks := content.Blocks.OfType("spec")
		switch len(specBlocks) {
		case 0:
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected component spec",
				Detail:   `A component block must have either a spec attribute or a nested spec block, but neither was found.`,
				Subject:  block.DefRange().Ptr(),
			})
		case 1:
			c.spec = string(hclgen.FormatBlock(specBlocks[0]))
		default:
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected at most one spec block",
				Detail:   fmt.Sprintf("Only one spec block may appear in a component, but found %d", len(specBlocks)),
				Subject:  specBlocks[1].DefRange.Ptr(),
			})
		}
	} else {
		c.spec = m.evalString(specAttr.Expr)
	}

	depsAttr := content.Attributes["depends_on"]
	if depsAttr != nil {
		depsExpr := depsAttr.Expr
		tup, ok := depsExpr.(*hclsyntax.TupleConsExpr)
		if !ok {
			m.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected array of strings",
				Detail:   fmt.Sprintf("Expected literal array of strings, got %T", depsExpr),
				Subject:  depsExpr.Range().Ptr(),
			})
		}
		c.dependsOn = make([]string, 0, len(tup.Exprs))
		for _, elem := range tup.Exprs {
			dep, diag := parseLiteralString(elem)
			if diag != nil {
				m.appendDiags(diag)
				continue
			}
			c.dependsOn = append(c.dependsOn, dep)
		}
	}

	return c
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) Type() string {
	return c.typ
}

func (c *Component) Spec() string {
	return c.spec
}

func (c *Component) DependsOn() []string {
	return c.dependsOn
}

func expandComponent(block *hclsyntax.Block) (*hclsyntax.Block, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	body := block.Body
	var encodefunc string
	switch block.Type {
	case "process":
		encodefunc = "jsonencode"
	case "container", "volume", "network":
		encodefunc = "yamlencode"
	default:
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported component type",
			Detail:   fmt.Sprintf(`The component type %q is not recognized.`, block.Type),
			Subject:  block.DefRange().Ptr(),
		})
		return nil, diags
	}
	if len(body.Blocks) > 0 {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected block",
			Detail:   fmt.Sprintf(`Unexpected block in %q component.`, block.Type),
			Subject:  body.Blocks[0].DefRange().Ptr(),
		})
	}
	attrs := body.Attributes
	specItems := make([]hclsyntax.ObjectConsItem, 0, len(attrs))
	for _, attr := range attrs {
		specItems = append(specItems, hclsyntax.ObjectConsItem{
			KeyExpr:   hclgen.NewObjStringKey(attr.Name, attr.Range()),
			ValueExpr: attr.Expr,
		})
	}
	// sort.Sort(specItemsSorter{specItems}) // XXX sort specItems by attr range?
	// XXX search for "_" blocks with depends_on, etc. and other meta properties.
	return &hclsyntax.Block{
		Type:   "component",
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
				"spec": &hclsyntax.Attribute{
					Name: "spec",
					Expr: &hclsyntax.FunctionCallExpr{
						Name: encodefunc,
						Args: []hclsyntax.Expression{
							&hclsyntax.ObjectConsExpr{
								Items:     specItems,
								SrcRange:  body.SrcRange,
								OpenRange: block.OpenBraceRange,
							},
						},
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
