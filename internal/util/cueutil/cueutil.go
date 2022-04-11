package cueutil

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
)

func EncodeValue(x any) cue.Value {
	cc := cuecontext.New()
	return cc.Encode(x)
}

func NodeToBytes(n ast.Node) ([]byte, error) {
	return format.Node(n)
}

func NodeToString(n ast.Node) (string, error) {
	bs, err := NodeToBytes(n)
	return string(bs), err
}

func ValueToBytes(v cue.Value, opts ...cue.Option) ([]byte, error) {
	node := v.Syntax(opts...)
	expr := node.(ast.Expr)
	file, err := astutil.ToFile(expr)
	if err != nil {
		return nil, err
	}
	return NodeToBytes(file)
}

func ValueToString(v cue.Value, opts ...cue.Option) (string, error) {
	bs, err := ValueToBytes(v, opts...)
	return string(bs), err
}

func MustValueToBytes(v cue.Value, opts ...cue.Option) []byte {
	bs, err := ValueToBytes(v, opts...)
	if err != nil {
		panic(err)
	}
	return bs
}

func MustValueToString(v cue.Value, opts ...cue.Option) string {
	s, err := ValueToString(v, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

func FormatBytes(bs []byte) ([]byte, error) {
	f, err := parser.ParseFile("", bs, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return NodeToBytes(f)
}

func FormatString(s string) (string, error) {
	f, err := parser.ParseFile("", s, parser.ParseComments)
	if err != nil {
		return "", err
	}
	return NodeToString(f)
}
