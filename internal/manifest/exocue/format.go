package exocue

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
)

func FormatBytes(v interface{}) ([]byte, error) {
	var node ast.Node
	switch v := v.(type) {
	case cue.Value:
		node = v.Syntax()
	case ast.Node:
		node = v
	default:
		panic(fmt.Errorf("cannot format %T", v))
	}
	return format.Node(node)
}

func FormatString(v interface{}) (string, error) {
	bs, err := FormatBytes(v)
	return string(bs), err
}

func StructToFile(v cue.Value) *ast.File {
	var decls []ast.Decl
	switch x := v.Syntax().(type) {
	case *ast.StructLit:
		decls = x.Elts
	case *ast.BottomLit:
		decls = []ast.Decl{x}
	default:
		panic(fmt.Errorf("cannot convert %T to file", x))
	}
	return &ast.File{
		Decls: decls,
	}
}

func dumpValue(v cue.Value) {
	bs, err := FormatBytes(v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bs)
}
