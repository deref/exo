package cueutil

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
)

func Final(v cue.Value) ast.Node {
	return StructToFile(v.Syntax(cue.Final()))
}

func EncodeValue(x any) cue.Value {
	cc := cuecontext.New()
	return cc.Encode(x)
}

func StructToFile(v any) *ast.File {
	var decls []ast.Decl
	switch v := v.(type) {
	case cue.Value:
		return StructToFile(v.Syntax())
	case *ast.StructLit:
		decls = v.Elts
	case *ast.BottomLit:
		decls = []ast.Decl{v}
	default:
		panic(fmt.Errorf("cannot convert %T to file", v))
	}
	return &ast.File{
		Decls: decls,
	}
}
