package codegen

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/deref/exo/inflect"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Root struct {
	Package string
	Module
}

type Module struct {
	Interfaces []Interface `hcl:"interface,block"`
	Structs    []Struct    `hcl:"struct,block"`
}

type Interface struct {
	Name    string   `hcl:"name,label"`
	Doc     *string  `hcl:"doc"`
	Methods []Method `hcl:"method,block"`
}

type Method struct {
	Name   string  `hcl:"name,label"`
	Doc    *string `hcl:"doc"`
	Input  []Field `hcl:"input,block"`
	Output []Field `hcl:"output,block"`
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

func ParseFile(filename string) (*Module, error) {
	var module Module
	if err := hclsimple.DecodeFile(filename, nil, &module); err != nil {
		return nil, err
	}
	err := Validate(module)
	return &module, err
}

func Validate(mod Module) error {
	// TODO: Validate no duplicate names.
	// TODO: All type references.
	return nil
}

func Generate(root *Root) ([]byte, error) {
	tmpl := template.Must(
		template.New("module").
			Funcs(map[string]interface{}{
				"tick":   func() string { return "`" },
				"pascal": inflect.KebabToPascal,
				"camel":  inflect.KebabToCamel,
			}).
			Parse(moduleTemplate),
	)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, root); err != nil {
		return nil, err
	}
	bs := buf.Bytes()
	formatted, err := format.Source(bs)
	if err != nil {
		return bs, nil
	}
	return formatted, nil

}

var moduleTemplate = `
// Generated file. DO NOT EDIT.

package {{.Package}}

{{- define "doc" -}}
{{if .Doc}}// {{.Doc}}
{{end}}{{end}}

{{- define "fields" -}}
{{- range . }}
{{template "doc" . -}}
	{{.Name|pascal}} {{.Type}} {{tick}}json:"{{.Name|camel}}"{{tick}}
{{- end}}{{end}}

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

{{range .Interfaces}}
{{template "doc" . -}}
type {{.Name|pascal}} interface {
{{- range .Methods}}
{{template "doc" . -}}
	{{.Name|pascal}}(context.Context, *{{.Name|pascal}}Input) (*{{.Name|pascal}}Output, error)
{{- end}}
}

{{range .Methods}}
type {{.Name|pascal}}Input struct {
{{template "fields" .Input}}
}

type {{.Name|pascal}}Output struct {
{{template "fields" .Output}}
}
{{end}}

func New{{.Name|pascal}}Mux(prefix string, iface {{.Name|pascal}}) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	Build{{.Name|pascal}}Mux(b, iface)
	return b.Mux()
}

func Build{{.Name|pascal}}Mux(b *josh.MuxBuilder, iface {{.Name|pascal}}) {
{{- range .Methods}}
	b.AddMethod("{{.Name}}", iface.{{.Name|pascal}})
{{- end}}
}
{{end}}

{{range .Structs}}
{{template "doc" . -}}
type {{.Name|pascal}} struct {
{{template "fields" .Fields}}
}
{{end}}
`
