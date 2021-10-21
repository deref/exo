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

type Converter struct {
	// ProjectName is used as a prefix for the resources created by this importer.
	ProjectName string
}

func (c *Converter) Convert(bs []byte) (*hcl.File, hcl.Diagnostics) {
	project, err := compose.Parse(bytes.NewBuffer(bs))
	if err != nil {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  err.Error(),
			},
		}
	}

	// TODO: Avoid mutating project during conversion.

	b := exohcl.NewBuilder(bs)
	var diags hcl.Diagnostics

	// Since containers reference networks and volumes by their docker-compose name, but the
	// Docker components will have a namespaced name, so we need to keep track of which
	// volumes/components a service references.
	networkKeyToName := map[string]string{}
	volumeKeyToName := map[string]string{}

	for _, volume := range project.Volumes {
		name := exohcl.MangleName(volume.Key)
		if volume.Key != name {
			var subject *hcl.Range
			diags = append(diags, exohcl.NewRenameWarning(volume.Key, name, subject))
		}

		if volume.Name.Value == "" {
			volume.Name = compose.MakeString(c.prefixedName(volume.Key, ""))
		}
		volumeKeyToName[volume.Key] = volume.Name.Value

		b.AddComponentBlock(makeComponentBlock("volume", name, volume, nil))
	}

	// Set up networks.
	hasDefaultNetwork := false
	for _, network := range project.Networks {
		if network.Key == "default" {
			hasDefaultNetwork = true
		}
		name := exohcl.MangleName(network.Key)
		if network.Key != name {
			var subject *hcl.Range
			diags = append(diags, exohcl.NewRenameWarning(network.Key, name, subject))
		}

		// If `name` is specified in the network configuration (usually used in conjunction with `external: true`),
		// then we should honor that as the docker network name. Otherwise, we should set the name as
		// `<project_name>_<network_key>`.
		if network.Name.Value == "" {
			network.Name = compose.MakeString(c.prefixedName(network.Key, ""))
		}
		networkKeyToName[network.Key] = network.Name.Value

		if network.Driver.Value == "" {
			network.Driver = compose.MakeString("bridge")
		}

		b.AddComponentBlock(makeComponentBlock("network", name, network, nil))
	}
	// TODO: Docker Compose only creates the default network if there is at least 1 service that does not
	// specify a network. We should do the same.
	if !hasDefaultNetwork {
		key := "default"
		name := c.prefixedName(key, "")
		networkKeyToName[key] = name

		b.AddComponentBlock(makeComponentBlock("network", key, map[string]string{
			"name":   name,
			"driver": "bridge",
		}, nil))
	}

	for _, service := range project.Services {
		name := exohcl.MangleName(service.Key)
		if service.Key != name {
			var subject *hcl.Range
			diags = append(diags, exohcl.NewRenameWarning(service.Key, name, subject))
		}
		var dependsOn []string

		if service.ContainerName.Value == "" {
			// The generated container name intentionally matches the container name generated by Docker Compose
			// when the scale is set at 1. When we address scaling containers, this will need to be updated to
			// use a different suffix for each container.
			service.ContainerName = compose.MakeString(c.prefixedName(service.Key, "1"))
		}

		for _, item := range service.Labels.Items {
			if strings.HasPrefix(item.Key, "com.docker.compose") {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("service may not specify labels with prefix \"com.docker.compose\", but %q specified %q", service.Key, item.Key),
				})
				return nil, diags
			}
		}
		// TODO: It probably makes more sense for these labels to be omitted here,
		// but preserved if adopting an existing Docker resource. As it is now,
		// they show up in converted manifests, which doesn't make much sense.
		service.Labels.Items = append(service.Labels.Items, compose.DictionaryItem{
			Key:   "com.docker.compose.project",
			Value: c.ProjectName,
		}, compose.DictionaryItem{
			Key:   "com.docker.compose.service",
			Value: service.Key,
		})

		// Map the docker-compose network name to the name of the docker network that is created.
		defaultNetworkName := networkKeyToName["default"]
		if len(service.Networks.Items) > 0 {
			mappedNetworks := make([]compose.ServiceNetwork, len(service.Networks.Items))
			for i, network := range service.Networks.Items {
				_, ok := networkKeyToName[network.Key]
				if !ok {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  fmt.Sprintf("unknown network: %q", network.Key),
					})
					continue
				}
				mappedNetworks[i] = network
				dependsOn = append(dependsOn, network.Key)
			}
			service.Networks.Items = mappedNetworks
		} else {
			service.Networks.Items = []compose.ServiceNetwork{
				{
					Key:       defaultNetworkName,
					ShortForm: compose.MakeString(defaultNetworkName),
				},
			}
			dependsOn = append(dependsOn, "default")
		}

		if len(service.Volumes) > 0 {
			for i, volumeMount := range service.Volumes {
				if volumeMount.Type.Value != "volume" {
					continue
				}
				if volumeName, ok := volumeKeyToName[volumeMount.Source.Value]; ok {
					originalName := volumeMount.Source.Value
					service.Volumes[i].Source = compose.MakeString(volumeName)
					dependsOn = append(dependsOn, originalName)
				}
				// If the volume was not listed in the top-level "volumes" section, then the docker engine
				// will create a new volume that will not be namespaced by the Compose project name.
			}
		}

		for _, dependency := range service.DependsOn.Items {
			condition := dependency.Condition
			if condition.Value == "" {
				condition = compose.MakeString("service_started")
			}
			if condition.Value != "service_started" {
				var subject *hcl.Range
				diags = append(diags, exohcl.NewUnsupportedFeatureWarning(
					fmt.Sprintf("service condition %q", dependency.Service),
					"only service_started is currently supported",
					subject,
				))
			}
			dependsOn = append(dependsOn, exohcl.MangleName(dependency.Service.Value))
		}

		for idx, link := range service.Links {
			var linkService, linkAlias string
			parts := strings.Split(link.Value, ":")
			switch len(parts) {
			case 1:
				linkService = parts[0]
				linkAlias = parts[0]
			case 2:
				linkService = parts[0]
				linkAlias = parts[1]
			default:
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("expected SERVICE or SERVICE:ALIAS for link, but got: %q", link),
				})
				return nil, diags
			}
			// NOTE [RESOLVING SERVICE CONTAINERS]:
			// There are several locations in a compose definition where a service may reference another service
			// by the compose name. We currently handle these situations by rewriting these locations to reference
			// a container named `<project>_<mangled_service_name>_1` with the assumption that a container will
			// be created by that name. However, this will break when the referenced service specifies a non-default
			// container name. Additionally, we may want to handle cases where a service is scaled past a single
			// container.
			// Some of these values could/should be resolved at runtime, and we should do it when we have the entire
			// project graph available.

			// See https://github.com/docker/compose/blob/v2.0.0-rc.3/compose/service.py#L836 for how compose configures
			// links.
			mangledServiceName := exohcl.MangleName(linkService)
			containerName := c.prefixedName(mangledServiceName, "1")
			service.Links[idx] = compose.Link{
				String:  compose.MakeString(fmt.Sprintf("%s:%s", containerName, linkAlias)),
				Service: containerName,
				Alias:   linkAlias,
			}
			dependsOn = append(dependsOn, mangledServiceName)
		}

		b.AddComponentBlock(makeComponentBlock("container", name, service, dependsOn))
	}

	return b.Build(), diags
}

func (c *Converter) prefixedName(name string, suffix string) string {
	var out strings.Builder
	out.WriteString(c.ProjectName)
	out.WriteByte('_')
	out.WriteString(name)
	if suffix != "" {
		out.WriteByte('_')
		out.WriteString(suffix)
	}

	return out.String()
}

func makeComponentBlock(typ string, name string, spec interface{}, dependsOn []string) *hclgen.Block {
	obj := yamlToHCL(spec).(*hclsyntax.ObjectConsExpr)
	attrs := make([]*hclsyntax.Attribute, len(obj.Items))
	for i, item := range obj.Items {
		key := item.KeyExpr.(*hclsyntax.ObjectConsKeyExpr)
		name := hcl.ExprAsKeyword(key.Wrapped)
		if name == "" {
			fmt.Printf("%#v\n", key.Wrapped.(*hclsyntax.TemplateExpr).Parts[0])
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
	var blocks []*hclgen.Block
	if len(dependsOn) > 0 {
		blocks = append(blocks, &hclgen.Block{
			Type: "_",
			Body: &hclgen.Body{
				Attributes: []*hclsyntax.Attribute{
					{
						Name: "depends_on",
						Expr: yamlToHCL(dependsOn),
					},
				},
			},
		})
	}
	return &hclgen.Block{
		Type:   typ,
		Labels: []string{name},
		Body: &hclgen.Body{
			Attributes: attrs,
			Blocks:     blocks,
		},
	}
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
