package exocue

import (
	_ "embed"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
	"github.com/deref/exo/internal/util/cacheutil"
)

//go:embed manifest.cue
var manifestSchemaSource string
var manifestSchemaFile = cacheutil.NewLazy(func() *ast.File {
	return mustParseFile("manifest.cue", manifestSchemaSource)
})

//go:embed resolution.cue
var resolutionSchemaSource string
var resolutionSchemaFile = cacheutil.NewLazy(func() *ast.File {
	return mustParseFile("resolution.cue", resolutionSchemaSource)
})

func mustParseFile(name string, content string) *ast.File {
	f, err := parser.ParseFile(name, content)
	if err != nil {
		panic(err)
	}
	return f
}

type Builder struct {
	decls []ast.Decl
}

func NewBuilder() *Builder {
	b := &Builder{}
	b.decls = append(b.decls, manifestSchemaFile.Force().Decls...)
	b.decls = append(b.decls, resolutionSchemaFile.Force().Decls...)
	return b
}

func declsToStruct(decls []ast.Decl) *ast.StructLit {
	return &ast.StructLit{
		Lbrace: token.NoSpace.Pos(),
		Elts:   decls,
	}
}

func (b *Builder) addDecl(path []string, decl ast.Decl) {
	for i := len(path) - 1; i >= 0; i-- {
		decl = ast.NewStruct(path[i], decl)
	}
	b.decls = append(b.decls, decl)
}

func (b *Builder) SetStack(id string, name string) {
	b.addDecl([]string{"$stack"}, ast.NewStruct(
		"id", id,
		"name", name,
	))
}

func (b *Builder) SetCluster(id string, name string, environment map[string]string) {
	envElems := make([]any, 0, len(environment)*2)
	for k, v := range environment {
		envElems = append(envElems, k, ast.NewString(v))
	}
	b.addDecl([]string{"$cluster"}, ast.NewStruct(
		"id", ast.NewString(id),
		"name", ast.NewString(name),
		"environment", ast.NewStruct(envElems...),
	))
}

func (b *Builder) AddComponent(id string, name string, typ string, spec cue.Value, parentID *string) {
	component := ast.NewStruct(
		"id", ast.NewString(id),
		"name", ast.NewString(name),
		"type", ast.NewString(typ),
		"spec", spec.Syntax(),
	)
	var decl ast.Expr
	switch typ {
	case "daemon":
		decl = newAnd(ast.NewIdent("#Daemon"), component)
	case "process":
		decl = newAnd(ast.NewIdent("#Process"), component)
	default:
		panic(fmt.Errorf("TODO: type lookup schema voodoo. typ=%q", typ))
	}
	b.addDecl([]string{"$components", id}, decl)
	sel := ast.NewSel(ast.NewIdent("$components"), id)
	b.addDecl([]string{"$components", id}, sel)
	if parentID == nil {
		b.addDecl([]string{"$stack", "components", name}, sel)
	} else {
		b.addDecl([]string{"$components", *parentID, "children", name}, sel)
	}
}

func (b *Builder) AddResource(id string, typ string, iri *string, componentID *string) {
	fields := []any{
		"id", ast.NewString(id),
		"type", ast.NewString(typ),
	}
	if iri != nil {
		fields = append(fields, "iri", ast.NewString(*iri))
	}
	resource := ast.NewStruct(fields...)
	b.addDecl([]string{"$resources", id}, resource)
	sel := ast.NewSel(ast.NewIdent("$resources"), id)
	if componentID == nil {
		b.addDecl([]string{"$stack", "detachedResources", id}, sel)
	} else {
		b.addDecl([]string{"$components", *componentID, "resources", id}, sel)
	}
}

func newAnd(xs ...ast.Expr) ast.Expr {
	return ast.NewBinExpr(token.AND, xs...)
}

func (b *Builder) Build() (Configuration, error) {
	cc := cuecontext.New()
	x := cc.BuildExpr(declsToStruct(b.decls))
	return Configuration(x), x.Validate()
}
