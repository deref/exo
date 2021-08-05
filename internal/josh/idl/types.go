package idl

type Package struct {
	Path string
	Unit
}

type Unit struct {
	Interfaces []Interface `hcl:"interface,block"`
	Structs    []Struct    `hcl:"struct,block"`
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
