package exohcl

import (
	"fmt"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type ComponentSet struct {
	// Analysis inputs.
	Blocks hcl.Blocks

	// Analysis outputs.
	Components []*Component
}

func NewComponentSet(m *Manifest) *ComponentSet {
	return &ComponentSet{
		Blocks: m.Components,
	}
}

func (cs *ComponentSet) Analyze(ctx *AnalysisContext) {
	if len(cs.Blocks) > 1 {
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Expected at most one components block",
			Detail:   fmt.Sprintf("Only one components block may appear in a manifest, but found %d", len(cs.Blocks)),
			Subject:  cs.Blocks[1].DefRange.Ptr(),
		})
	}

	for _, block := range cs.Blocks {
		if len(block.Labels) > 0 {
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected label on components block",
				Detail:   fmt.Sprintf("A components block expects no labels, but has %d", len(block.Labels)),
				Subject:  &block.LabelRanges[0],
			})
		}
		body, ok := block.Body.(*hclsyntax.Body)
		if !ok {
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Malformed components block",
				Detail:   fmt.Sprintf("Expected components block to be an *hclsyntax.Body, but got %T", block.Body),
				Subject:  &block.DefRange,
			})
			continue
		}
		if len(body.Attributes) > 0 {
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected attributes in components block",
				Detail:   fmt.Sprintf("A components block expects no attributes, but has %d", len(body.Attributes)),
				Subject:  body.Attributes.Range().Ptr(),
			})
		}
		for _, componentBlock := range body.Blocks {
			component := NewComponent(componentBlock)
			component.Analyze(ctx)
			cs.Components = append(cs.Components, component)
		}
	}
}

type Component struct {
	Source *hclsyntax.Block

	Expansion *hclsyntax.Block
	Type      string
	Name      string
	Spec      string
	DependsOn []string
}

func NewComponent(block *hclsyntax.Block) *Component {
	return &Component{
		Source: block,
	}
}

func (c *Component) Analyze(ctx *AnalysisContext) {
	c.Expansion = expandComponent(ctx, c.Source)
	block := c.Expansion
	if block == nil {
		return
	}

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
	ctx.AppendDiags(diags...)

	switch len(block.Labels) {
	case 0:
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected component name",
			Detail:   `A component block must have exactly one label, which is the name of the component.`,
			Subject:  block.DefRange().Ptr(),
		})
	case 1:
		c.Name = block.Labels[0]
	default:
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected label.",
			Detail:   `A component block must have exactly one label, which is the name of the component.`,
			Subject:  block.LabelRanges[1].Ptr(),
		})
	}

	typeAttr := content.Attributes["type"]
	if typeAttr != nil {
		var diag *hcl.Diagnostic
		c.Type, diag = parseLiteralString(typeAttr.Expr)
		if diag != nil {
			ctx.AppendDiags(diag)
		}
	}

	specAttr := content.Attributes["spec"]
	if specAttr == nil {
		specBlocks := content.Blocks.OfType("spec")
		switch len(specBlocks) {
		case 0:
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected component spec",
				Detail:   `A component block must have either a spec attribute or a nested spec block, but neither was found.`,
				Subject:  block.DefRange().Ptr(),
			})
		case 1:
			c.Spec = string(hclgen.FormatBlock(specBlocks[0]))
		default:
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected at most one spec block",
				Detail:   fmt.Sprintf("Only one spec block may appear in a component, but found %d", len(specBlocks)),
				Subject:  specBlocks[1].DefRange.Ptr(),
			})
		}
	} else {
		c.Spec, _ = AnalyzeString(ctx, specAttr.Expr)
	}

	depsAttr := content.Attributes["depends_on"]
	if depsAttr != nil {
		depsExpr := depsAttr.Expr
		tup, ok := depsExpr.(*hclsyntax.TupleConsExpr)
		if !ok {
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected array of strings",
				Detail:   fmt.Sprintf("Expected literal array of strings, got %T", depsExpr),
				Subject:  depsExpr.Range().Ptr(),
			})
		}
		c.DependsOn = make([]string, 0, len(tup.Exprs))
		for _, elem := range tup.Exprs {
			dep, diag := parseLiteralString(elem)
			if diag != nil {
				ctx.AppendDiags(diag)
				continue
			}
			c.DependsOn = append(c.DependsOn, dep)
		}
	}
}

func expandComponent(ctx *AnalysisContext, block *hclsyntax.Block) *hclsyntax.Block {
	body := block.Body
	var encodefunc string
	switch block.Type {
	case "component":
		return block
	case "process":
		encodefunc = "jsonencode"
	case "container", "volume", "network", "apigateway":
		encodefunc = "yamlencode"
	default:
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported component type",
			Detail:   fmt.Sprintf(`The component type %q is not recognized.`, block.Type),
			Subject:  block.DefRange().Ptr(),
		})
		return nil
	}
	for _, subblock := range body.Blocks {
		switch subblock.Type {
		case "_":
			// TODO: copy content of "meta" blocks in to the expanded output.
		default:
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected block",
				Detail:   fmt.Sprintf(`Unexpected %q block in %q component.`, subblock.Type, block.Type),
				Subject:  body.Blocks[0].DefRange().Ptr(),
			})
		}
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
	}
}
