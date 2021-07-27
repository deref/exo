package codegen

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/deref/exo/inflect"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Package struct {
	Path string
	Unit
}

type Unit struct {
	Interfaces  []Interface  `hcl:"interface,block"`
	Structs     []Struct     `hcl:"struct,block"`
	Controllers []Controller `hcl:"controller,block"`
}

type Interface struct {
	Name    string   `hcl:"name,label"`
	Doc     *string  `hcl:"doc"`
	Extends []string `hcl:"extends,optional"`
	Methods []Method `hcl:"method,block"`
}

type Method struct {
	Name    string  `hcl:"name,label"`
	Doc     *string `hcl:"doc"`
	Inputs  []Field `hcl:"input,block"`
	Outputs []Field `hcl:"output,block"`
}

type Struct struct {
	Name   string  `hcl:"name,label"`
	Doc    *string `hcl:"doc"`
	Fields []Field `hcl:"field,block"`
}

type Field struct {
	Name     string  `hcl:"name,label"`
	Doc      *string `hcl:"doc"`
	Type     string  `hcl:"type,label"`
	Required *bool   `hcl:"required"`
	Nullable *bool   `hcl:"nullable"`
}

type Controller struct {
	Name    string   `hcl:"name,label"`
	Doc     *string  `hcl:"doc"`
	Extends []string `hcl:"extends,optional"`
	Methods []Method `hcl:"method,block"`
}

func LoadFile(filename string) (*Unit, error) {
	unit, err := ParseFile(filename)
	if err != nil {
		return nil, err
	}
	Elaborate(unit)
	err = Validate(unit)
	return unit, err
}

func ParseFile(filename string) (*Unit, error) {
	var unit Unit
	if err := hclsimple.DecodeFile(filename, nil, &unit); err != nil {
		return nil, err
	}
	return &unit, nil
}

func Elaborate(unit *Unit) {
	for _, controller := range unit.Controllers {
		methods := make([]Method, len(controller.Methods))
		for methodIndex, method := range controller.Methods {
			shiftInputs := 3
			inputs := make([]Field, len(method.Inputs)+shiftInputs)
			inputs[0] = Field{
				Name: "id",
				Type: "string",
			}
			inputs[1] = Field{
				Name: "spec",
				Type: "string",
			}
			inputs[2] = Field{
				Name: "state",
				Type: "string",
			}
			for inputIndex, input := range method.Inputs {
				inputs[shiftInputs+inputIndex] = input
			}

			shiftOutputs := 1
			outputs := make([]Field, len(method.Outputs)+shiftOutputs)
			outputs[0] = Field{
				Name: "state",
				Type: "string",
			}
			for outputIndex, output := range method.Outputs {
				inputs[shiftOutputs+outputIndex] = output
			}

			methods[methodIndex] = Method{
				Name:    method.Name,
				Doc:     method.Doc,
				Inputs:  inputs,
				Outputs: outputs,
			}
		}
		unit.Interfaces = append(unit.Interfaces, Interface{
			Name:    controller.Name,
			Doc:     controller.Doc,
			Extends: append([]string{}, controller.Extends...),
			Methods: methods,
		})
	}
}

func Validate(unit *Unit) error {
	// TODO: Validate no duplicate names.
	// TODO: All type references.
	return nil
}

func GenerateAPI(pkg *Package) ([]byte, error) {
	return generateGo(apiTemplate, pkg)
}

func GenerateClient(pkg *Package) ([]byte, error) {
	return generateGo(clientTemplate, pkg)
}

func generateGo(t string, pkg *Package) ([]byte, error) {
	tmpl := template.Must(
		template.New("package").
			Funcs(templateFuncs).
			Parse(t),
	)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, pkg); err != nil {
		return nil, err
	}
	bs := buf.Bytes()
	formatted, err := format.Source(bs)
	if err != nil {
		return bs, nil
	}
	return formatted, nil
}

var templateFuncs = map[string]interface{}{
	"tick":   func() string { return "`" },
	"public": inflect.KebabToPublic,
	"js":     inflect.KebabToJSVar,
}

var apiTemplate = `
// Generated file. DO NOT EDIT.

package api

{{- define "doc" -}}
{{if .Doc}}// {{.Doc}}
{{end}}{{end}}

{{- define "fields" -}}
{{- range . }}
{{template "doc" . -}}
	{{.Name|public}} {{.Type}} {{tick}}json:"{{.Name|js}}"{{tick}}
{{- end}}{{end}}

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

{{range .Interfaces}}
{{template "doc" . -}}
type {{.Name|public}} interface {
{{- range .Extends}}
	{{.|public}}
{{- end}}
{{- range .Methods}}
{{template "doc" . -}}
	{{.Name|public}}(context.Context, *{{.Name|public}}Input) (*{{.Name|public}}Output, error)
{{- end}}
}

{{range .Methods}}
type {{.Name|public}}Input struct {
{{template "fields" .Inputs}}
}

type {{.Name|public}}Output struct {
{{template "fields" .Outputs}}
}
{{end}}

func Build{{.Name|public}}Mux(b *josh.MuxBuilder, factory func(req *http.Request) {{.Name|public}}) {
{{- range .Extends}}
	Build{{.|public}}Mux(b, func(req *http.Request) {{.|public}} {
		return factory(req)
	})
{{- end }}
{{- range .Methods}}
	b.AddMethod("{{.Name}}", func (req *http.Request) interface{} {
		return factory(req).{{.Name|public}}
	})
{{- end}}
}
{{end}}

{{range .Structs}}
{{template "doc" . -}}
type {{.Name|public}} struct {
{{template "fields" .Fields}}
}
{{end}}
`

var clientTemplate = `
// Generated file. DO NOT EDIT.

package client

import (
	"context"

	josh "github.com/deref/exo/josh/client"
	"github.com/deref/{{.Path}}/api"
)

{{range $_, $iface := .Interfaces}}
type {{.Name|public}} struct {
	client *josh.Client
}

var _ api.{{.Name|public}} = (*{{.Name|public}})(nil)

func Get{{.Name|public}}(client *josh.Client) *{{.Name|public}} {
	return &{{.Name|public}}{
		client: client,
	}
}

{{range .Methods}}
func (c *{{$iface.Name|public}}) {{.Name|public}}(ctx context.Context, input *api.{{.Name|public}}Input) (output *api.{{.Name|public}}Output, err error) {
	err = c.client.Invoke(ctx, "{{.Name}}", input, &output)
	return
}
{{end}}
{{end}}
`
