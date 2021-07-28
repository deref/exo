package model

type Struct struct {
	pkg    *Package
	name   string
	doc    string
	fields []*Field
}

func newStruct(pkg *Package, name string) *Struct {
	return &Struct{
		pkg:  pkg,
		name: name,
	}
}

func (strct *Struct) Name() string {
	return strct.name
}

func (strct *Struct) SetDoc(value string) {
	strct.doc = value
}

func (strct *Struct) Doc() string {
	return strct.doc
}

func (strct *Struct) Fields() []*Field {
	fields := make([]*Field, len(strct.fields))
	copy(fields, strct.fields)
	return fields
}

func (strct *Struct) AddField(cfg FieldConfig) *Field {
	field := &Field{cfg: cfg}
	strct.fields = append(strct.fields, field)
	return field
}
