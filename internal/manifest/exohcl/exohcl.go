package exohcl

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	ctyyaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var Latest = FormatVersion{
	Major: 1,
	Minor: 0,
}

const Starter = `# See https://docs.deref.io/exo for details.

exo = "1.0"
`

type Manifest struct {
	FormatVersion FormatVersion
	Components    []Component
}

type FormatVersion struct {
	Major int
	Minor int
}

func (ver FormatVersion) String() string {
	return fmt.Sprintf("%d.%d", ver.Major, ver.Minor)
}

type Component struct {
	Name      string
	Type      string
	Spec      string
	DependsOn []string
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

type Loader struct {
	Filename string
	diags    hcl.Diagnostics
	evalCtx  *hcl.EvalContext
}

func (l *Loader) LoadBytes(bs []byte) (*Manifest, error) {
	var file *hcl.File
	file, l.diags = hclsyntax.ParseConfig(bs, l.Filename, hcl.InitialPos)
	manifest := l.loadHCL(file)
	return &manifest, l.diags
}

func (l *Loader) appendDiags(diags ...*hcl.Diagnostic) {
	l.diags = append(l.diags, diags...)
}

func (l *Loader) LoadHCL(file *hcl.File) (*Manifest, error) {
	m := l.loadHCL(file)
	return &m, l.diags
}

func (l *Loader) loadHCL(file *hcl.File) (manifest Manifest) {
	l.evalCtx = &hcl.EvalContext{
		Functions: map[string]function.Function{
			"jsonencode": stdlib.JSONEncodeFunc,
			"yamlencode": ctyyaml.YAMLEncodeFunc,
		},
	}

	if file == nil {
		return
	}
	content, diags := file.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "exo", Required: true},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "components"},
		},
	})
	l.appendDiags(diags...)
	if content == nil {
		return
	}

	versionAttr := content.Attributes["exo"]
	if versionAttr != nil {
		manifest.FormatVersion = l.parseFormatVersion(versionAttr)
	}

	containers := content.Blocks.OfType("components")
	switch len(containers) {
	case 0:
		manifest.Components = []Component{}
	case 1:
		manifest.Components = l.parseComponents(containers[0])
	default:
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected at most one components block",
			Detail:   fmt.Sprintf("Only one components block may appear in a manifest, but found %d", len(containers)),
			Subject:  containers[1].DefRange.Ptr(),
		})
	}
	return
}

func (l *Loader) parseFormatVersion(attr *hcl.Attribute) (version FormatVersion) {
	s, diag := l.parseLiteralString(attr.Expr)
	if diag != nil {
		l.appendDiags(diag)
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
		l.appendDiags(&hcl.Diagnostic{
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
			l.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported exo format minor version",
				Detail:   fmt.Sprintf(`Unsupported exo format minor version. The maximum supported "1.x" version is %q.`, Latest),
				Subject:  attr.Expr.Range().Ptr(),
				Context:  attr.Range.Ptr(),
			})
		}
	default:
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported exo format major version",
			Detail:   fmt.Sprintf(`Unsupported exo format major version. The latest version is %q.`, Latest),
			Subject:  attr.Expr.Range().Ptr(),
			Context:  attr.Range.Ptr(),
		})
	}
	return
}

func (l *Loader) parseComponents(block *hcl.Block) (components []Component) {
	if len(block.Labels) > 0 {
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected label on components block",
			Detail:   fmt.Sprintf("A components block expects no labels, but has %d", len(block.Labels)),
			Subject:  &block.LabelRanges[0],
		})
	}
	body, ok := block.Body.(*hclsyntax.Body)
	if !ok {
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Malformed components block",
			Detail:   fmt.Sprintf("Expected components block to be an *hclsyntax.Body, but got %T", block.Body),
			Subject:  &block.DefRange,
		})
		return
	}
	if len(body.Attributes) > 0 {
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected attributes in components block",
			Detail:   fmt.Sprintf("A components block expects no attributes, but has %d", len(body.Attributes)),
			Subject:  body.Attributes.Range().Ptr(),
		})
	}
	components = make([]Component, 0, len(body.Blocks))
	for _, block := range body.Blocks {
		component := l.parseComponent(block)
		if component != nil {
			components = append(components, *component)
		}
	}
	return
}

func (l *Loader) parseComponent(block *hclsyntax.Block) *Component {
	if block.Type != "component" {
		block = l.expandComponent(block)
		if block == nil {
			return nil
		}
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
	l.appendDiags(diags...)

	var component Component

	switch len(block.Labels) {
	case 0:
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected component name",
			Detail:   `A component block must have exactly one label, which is the name of the component.`,
			Subject:  block.DefRange().Ptr(),
		})
	case 1:
		component.Name = block.Labels[0]
	default:
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unexpected label.",
			Detail:   `A component block must have exactly one label, which is the name of the component.`,
			Subject:  block.LabelRanges[1].Ptr(),
		})
	}

	typeAttr := content.Attributes["type"]
	if typeAttr != nil {
		var diag *hcl.Diagnostic
		component.Type, diag = l.parseLiteralString(typeAttr.Expr)
		if diag != nil {
			l.appendDiags(diag)
		}
	}

	specAttr := content.Attributes["spec"]
	if specAttr == nil {
		specBlocks := content.Blocks.OfType("spec")
		switch len(specBlocks) {
		case 0:
			l.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected component spec",
				Detail:   `A component block must have either a spec attribute or a nested spec block, but neither was found.`,
				Subject:  block.DefRange().Ptr(),
			})
		case 1:
			component.Spec = formatBlock(specBlocks[0])
		default:
			l.appendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Expected at most one spec block",
				Detail:   fmt.Sprintf("Only one spec block may appear in a component, but found %d", len(specBlocks)),
				Subject:  specBlocks[1].DefRange.Ptr(),
			})
		}
	} else {
		component.Spec = l.evalString(specAttr.Expr)
	}

	depsAttr := content.Attributes["depends_on"]
	if depsAttr != nil {
		component.DependsOn = l.parseDependsOn(depsAttr.Expr)
	}

	return &component
}

func formatBlock(block *hcl.Block) string {
	panic("TODO: use hclwrite to format this block as text")
}

func (l *Loader) parseLiteralString(x hcl.Expression) (string, *hcl.Diagnostic) {
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

func (l *Loader) evalString(x hcl.Expression) string {
	v, diags := x.Value(l.evalCtx)
	l.appendDiags(diags...)
	if v.Type() != cty.String {
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected string",
			Detail:   fmt.Sprintf("Expected string, but got %s", v.Type().FriendlyName()),
			Subject:  x.Range().Ptr(),
		})
		return ""
	}
	return v.AsString()
}

func (l *Loader) parseDependsOn(x hcl.Expression) []string {
	tup, ok := x.(*hclsyntax.TupleConsExpr)
	if !ok {
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Expected array of strings",
			Detail:   fmt.Sprintf("Expected literal array of strings, got %T", x),
			Subject:  x.Range().Ptr(),
		})
	}
	res := make([]string, 0, len(tup.Exprs))
	for _, elem := range tup.Exprs {
		dep, diag := l.parseLiteralString(elem)
		if diag != nil {
			l.appendDiags(diag)
			continue
		}
		res = append(res, dep)
	}
	return res
}

func (l *Loader) expandComponent(block *hclsyntax.Block) *hclsyntax.Block {
	body := block.Body
	var encodefunc string
	switch block.Type {
	case "process":
		encodefunc = "jsonencode"
	case "container", "volume", "network":
		encodefunc = "yamlencode"
	default:
		l.appendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported component type",
			Detail:   fmt.Sprintf(`The component type %q is not recognized.`, block.Type),
			Subject:  block.DefRange().Ptr(),
		})
		return nil
	}
	if len(body.Blocks) > 0 {
		l.appendDiags(&hcl.Diagnostic{
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
			KeyExpr:   NewStringLiteral(attr.Name, attr.Range()),
			ValueExpr: attr.Expr,
		})
	}
	// sort.Sort(specItemsSorter{specItems}) // XXX sort specItems by attr range?
	return &hclsyntax.Block{
		Type:   "component",
		Labels: block.Labels,
		Body: &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{
				"type": &hclsyntax.Attribute{
					Name:        "type",
					Expr:        NewStringLiteral(block.Type, block.TypeRange),
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

func Generate(w io.Writer, manifest *Manifest) error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(manifest, f.Body())
	_, err := f.WriteTo(w)
	return err
}
