package compose

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"
)

type Converter struct{}

func (c *Converter) Convert(bs []byte) (*hcl.File, hcl.Diagnostics) {
	project, err := compose.Parse(bytes.NewBuffer(bs))
	if err != nil {
		// TODO: Get position information from YAML parser errors.
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  err.Error(),
			},
		}
	}

	b := exohcl.NewBuilder(bs)
	var diags hcl.Diagnostics

	// Compose has one namespace per "section" (ie type of component), but Exo
	// components live within a single namespace. Count these keys to detect
	// conflicts.
	keyCounts := make(map[string]int)
	addKey := func(key string) {
		keyCounts[key]++
	}
	for _, item := range project.Services {
		addKey(item.Key)
	}
	for _, item := range project.Networks {
		addKey(item.Key)
	}
	for _, item := range project.Volumes {
		addKey(item.Key)
	}
	for _, item := range project.Configs {
		addKey(item.Key)
	}
	for _, item := range project.Secrets {
		addKey(item.Key)
	}

	convertComponent := func(typ string, key string, spec interface{}) {
		name := exohcl.MangleName(key)
		// Detect conflicts between sections, report them against non-services.
		if typ != "service" && keyCounts[key] > 1 {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("%s name conflicts with service: %q", typ, key),
			})
		} else if key != name {
			// Renames are an error because intra-compose-file references can be
			// broken and we don't have a mechanism for applying renames to
			// references, such as in `services[*].depends_on` or
			// `service[*].volumes`.
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("invalid name: %q", key),
			})
		}

		obj := yamlToHCL(spec).(*hclsyntax.ObjectConsExpr)
		attrs := make([]*hclsyntax.Attribute, len(obj.Items))
		for i, item := range obj.Items {
			key := item.KeyExpr.(*hclsyntax.ObjectConsKeyExpr)
			name := hcl.ExprAsKeyword(key.Wrapped)
			if name == "" {
				panic("unexpected complex key")
			}
			val := item.ValueExpr
			attrs[i] = &hclsyntax.Attribute{
				Name:      name,
				NameRange: key.Range(),
				Expr:      val,
				SrcRange:  hcl.RangeBetween(key.Range(), val.Range()),
			}
		}
		block := &hclgen.Block{
			Type:   typ,
			Labels: []string{name},
			Body: &hclgen.Body{
				Attributes: attrs,
			},
		}
		b.AddComponentBlock(block)
	}

	for _, service := range project.Services {
		// We directly translate services to containers until we have a richer
		// notion of "service".
		convertComponent("container", service.Key, service)
	}
	for _, volume := range project.Volumes {
		convertComponent("volume", volume.Key, volume)
	}
	for _, network := range project.Networks {
		convertComponent("network", network.Key, network)
	}
	if len(project.Configs) > 0 {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  fmt.Sprintf(`compose "config" seciton not yet supported`),
		})
	}
	if len(project.Secrets) > 0 {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  fmt.Sprintf(`compose "secrets" seciton not yet supported`),
		})
	}

	return b.Build(), diags
}

func yamlToHCL(v interface{}) hclsyntax.Expression {
	if v == nil {
		return hclgen.NewNullLiteral(hcl.Range{})
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return hclgen.NewNullLiteral(hcl.Range{})
		}
		rv = rv.Elem()
	}
	v = rv.Interface()

	marshaller, ok := v.(yaml.Marshaler)
	if ok {
		v, err := marshaller.MarshalYAML()
		if err != nil {
			panic(err)
		}
		return yamlToHCL(v)
	}

	switch v := v.(type) {
	case string:
		return hclgen.NewStringLiteral(v, hcl.Range{})
	case int:
		return &hclsyntax.LiteralValueExpr{
			Val: cty.NumberIntVal(int64(v)),
		}
	case int16:
		return &hclsyntax.LiteralValueExpr{
			Val: cty.NumberIntVal(int64(v)),
		}
	case uint16:
		return &hclsyntax.LiteralValueExpr{
			Val: cty.NumberIntVal(int64(v)),
		}
	case int64:
		return &hclsyntax.LiteralValueExpr{
			Val: cty.NumberIntVal(v),
		}
	case float64:
		return &hclsyntax.LiteralValueExpr{
			Val: cty.NumberFloatVal(v),
		}
	case bool:
		return &hclsyntax.LiteralValueExpr{
			Val: cty.BoolVal(v),
		}
	case yaml.Node:
		switch v.Kind {
		case yaml.SequenceNode:
			elems := make([]hclsyntax.Expression, len(v.Content))
			for i, c := range v.Content {
				elems[i] = yamlToHCL(c)
			}
			return &hclsyntax.TupleConsExpr{
				Exprs: elems,
			}
		case yaml.MappingNode:
			n := len(v.Content) / 2
			items := make([]hclsyntax.ObjectConsItem, n)
			for i := 0; i < n; i++ {
				items[i] = hclsyntax.ObjectConsItem{
					KeyExpr:   yamlToHCLKey(v.Content[i*2+0]),
					ValueExpr: yamlToHCL(v.Content[i*2+1]),
				}
			}
			return &hclsyntax.ObjectConsExpr{
				Items: items,
			}

		case yaml.ScalarNode:
			if v.Tag != "" && v.Tag != "!!str" {
				panic(fmt.Errorf("unexpected yaml node tag: %q", v.Tag))
			}
			return yamlToHCL(v.Value)
		default:
			panic(fmt.Errorf("unexpected yaml node kind: %d", v.Kind))
		}
	default:
		switch rv.Kind() {
		case reflect.Struct:

			typ := rv.Type()

			numField := rv.NumField()
			var items []hclsyntax.ObjectConsItem
			for i := 0; i < numField; i++ {
				fld := typ.Field(i)
				tag := fld.Tag.Get("yaml")
				if tag == "" {
					panic(fmt.Errorf("struct %s field %s missing yaml tag", typ, fld.Name))
				}
				options := strings.Split(tag, ",")
				name := options[0]
				if name == "-" {
					continue
				}
				omitempty := false
				for _, option := range options[1:] {
					switch option {
					case "omitempty":
						omitempty = true
					default:
						panic(fmt.Errorf("unsupported yaml field tag option: %q", option))
					}
				}

				fldV := rv.Field(i)
				valueExpr := yamlToHCL(fldV.Interface())
				if omitempty && isZeroExpr(valueExpr) {
					continue
				}

				item := hclsyntax.ObjectConsItem{
					KeyExpr:   hclgen.NewObjStringKey(name, hcl.Range{}),
					ValueExpr: valueExpr,
				}
				items = append(items, item)
			}
			return &hclsyntax.ObjectConsExpr{
				Items: items,
			}

		case reflect.Slice:
			elems := make([]hclsyntax.Expression, rv.Len())
			for i := range elems {
				elems[i] = yamlToHCL(rv.Index(i).Interface())
			}
			return &hclsyntax.TupleConsExpr{
				Exprs: elems,
			}

		case reflect.Map:
			// XXX need a stable sort!
			items := make([]hclsyntax.ObjectConsItem, 0, rv.Len())
			iter := rv.MapRange()
			for iter.Next() {
				item := hclsyntax.ObjectConsItem{
					KeyExpr:   yamlToHCLKey(iter.Key().Interface()),
					ValueExpr: yamlToHCL(iter.Value().Interface()),
				}
				items = append(items, item)
			}
			return &hclsyntax.ObjectConsExpr{
				Items: items,
			}

		default:
			panic(fmt.Errorf("unexpected yaml type: %T", v))
		}
	}
}

func yamlToHCLKey(v interface{}) hclsyntax.Expression {
	x := yamlToHCL(v)
	if template, isTemplate := x.(*hclsyntax.TemplateExpr); isTemplate && len(template.Parts) == 1 {
		if lit, isLit := template.Parts[0].(*hclsyntax.LiteralValueExpr); isLit && lit.Val.Type() == cty.String {
			return hclgen.NewObjStringKey(lit.Val.AsString(), hcl.Range{})
		}
	}
	return &hclsyntax.ObjectConsKeyExpr{
		Wrapped: x,
	}
}

func isZeroExpr(x hcl.Expression) bool {
	v, err := x.Value(&hcl.EvalContext{})
	if err != nil {
		return false
	}
	typ := v.Type()
	switch {
	case v.Equals(cty.NilVal).True():
		return true
	case typ == cty.Number:
		return v.Equals(cty.Zero).True()
	case typ == cty.String:
		return v.AsString() == ""
	case typ == cty.Bool:
		return v.False()
	case typ.IsTupleType():
		return v.LengthInt() == 0
	case typ.IsCollectionType():
		return v.LengthInt() == 0
	case typ.IsObjectType():
		return v.LengthInt() == 0
	default:
		panic(fmt.Errorf("unexpected cty type: %s", typ))
	}
}
