package cueutil

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
)

func FormatBytes(v any) ([]byte, error) {
	var node ast.Node
	switch v := v.(type) {
	case cue.Value:
		node = v.Syntax()
	case ast.Node:
		node = v
	case string:
		return FormatBytes([]byte(v))
	case []byte:
		f, err := parser.ParseFile("", v, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		return FormatBytes(f)
	default:
		panic(fmt.Errorf("cannot format %T", v))
	}
	return format.Node(node)
}

func FormatString(v any) (string, error) {
	bs, err := FormatBytes(v)
	return string(bs), err
}

func MustFormatString(v any) string {
	s, err := FormatString(v)
	if err != nil {
		panic(err)
	}
	return s
}

func Dump(v any) {
	fmt.Println(MustFormatString(v))
}
