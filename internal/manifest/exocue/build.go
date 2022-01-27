package exocue

import (
	_ "embed"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/parser"
)

//go:embed schema.cue
var schema string

type Builder struct {
	decls []ast.Decl
}

func NewBuilder() *Builder {
	b := &Builder{}
	schema, err := parser.ParseFile("schema.cue", schema)
	if err != nil {
		panic(err)
	}
	b.addDecls(schema.Decls...)
	return b
}

func (b *Builder) addDecls(decls ...ast.Decl) {
	b.decls = append(b.decls, decls...)
}

func (b *Builder) AddManifest(s string) {
	manifest, err := parser.ParseFile("exo.cue", s, parser.ParseComments)
	if err != nil {
		panic(err) // XXX
	}
	b.addDecls(&ast.StructLit{
		Elts: []ast.Decl{
			&ast.Field{
				Label: ast.NewString("$stack"),
				Value: &ast.StructLit{
					Elts: manifest.Decls,
				},
			},
		},
	})
}

func (b *Builder) Build() *Configuration {
	cc := cuecontext.New()
	return &Configuration{
		v: cc.BuildExpr(&ast.StructLit{Elts: b.decls}),
	}
}
