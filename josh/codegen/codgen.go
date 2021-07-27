package codegen

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/deref/exo/inflect"
	"github.com/deref/exo/josh/model"
)

func GenerateAPI(pkg *model.Package) ([]byte, error) {
	return generateGo(apiTemplate, pkg)
}

func GenerateClient(pkg *model.Package) ([]byte, error) {
	return generateGo(clientTemplate, pkg)
}

func generateGo(t string, pkg *model.Package) ([]byte, error) {
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
