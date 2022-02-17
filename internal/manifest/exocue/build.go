package exocue

import (
	_ "embed"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
)

//go:embed schema.cue
var schema string

type Builder struct {
	decls []ast.Decl
}

func NewBuilder() *Builder {
	schema, err := parser.ParseFile("schema.cue", schema)
	if err != nil {
		panic(err)
	}
	return &Builder{
		decls: append([]ast.Decl{}, schema.Decls...),
	}
}

func declsToStruct(decls []ast.Decl) *ast.StructLit {
	return &ast.StructLit{
		Lbrace: token.NoSpace.Pos(),
		Elts:   decls,
	}
}

func parseFileAsStruct(fname, s string) *ast.StructLit {
	f, err := parser.ParseFile("exo.cue", s, parser.ParseComments)
	if err != nil {
		panic(err) // XXX
	}
	return declsToStruct(f.Decls)
}

func (b *Builder) addDecl(path []string, decl ast.Decl) {
	for i := len(path) - 1; i >= 0; i-- {
		decl = ast.NewStruct(path[i], decl)
	}
	b.decls = append(b.decls, decl)
}

func (b *Builder) AddManifest(s string) {
	manifest := parseFileAsStruct("exo.cue", s)
	b.addDecl([]string{"$stack"}, manifest)
}

func (b *Builder) AddComponent(id string, name string, typ string, spec string) {
	fname := fmt.Sprintf("components/%s.cue", name)
	specNode := parseFileAsStruct(fname, spec)
	res := ast.NewStruct(
		"id", ast.NewString(id),
		"type", ast.NewString(typ),
		"spec", specNode,
	)
	var decl ast.Expr
	switch typ {
	case "daemon":
		decl = newAnd(ast.NewIdent("$Daemon"), res)
	default:
		panic(fmt.Errorf("TODO: type lookup schema voodoo. typ=%q", typ))
	}
	b.addDecl([]string{"$stack", "components", name}, decl)
}

func (b *Builder) AddCluster(id string, name string, environment map[string]interface{}) {
	envElems := make([]interface{}, 0, len(environment)*2)
	for k, v := range environment {
		envElems = append(envElems, k, ast.NewString(v.(string)))
	}
	b.addDecl([]string{"$cluster"}, ast.NewStruct(
		"id", ast.NewString(id),
		"name", ast.NewString(name),
		"environment", ast.NewStruct(envElems...),
	))
}

func newAnd(xs ...ast.Expr) ast.Expr {
	return ast.NewBinExpr(token.AND, xs...)
}

func (b *Builder) build(key string) cue.Value {
	cc := cuecontext.New()
	cfg := cc.BuildExpr(declsToStruct(b.decls))
	return cfg.LookupPath(cue.MakePath(cue.Str(key)))
}

func (b *Builder) BuildStack() Stack {
	return Stack(b.build("$stack"))
}

func (b *Builder) BuildCluster() Cluster {
	return Cluster(b.build("$cluster"))
}
