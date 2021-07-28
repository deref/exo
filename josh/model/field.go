package model

type Field struct {
	cfg FieldConfig
}

type FieldConfig struct {
	Name     string
	Doc      string
	Type     string // TODO: Type type.
	Required bool
	Nullable bool
}

func (field *Field) Name() string {
	return field.cfg.Name
}

func (field *Field) Doc() string {
	return field.cfg.Doc
}

func (field *Field) Type() string {
	return field.cfg.Type
}

func (field *Field) Required() bool {
	return field.cfg.Required
}

func (field *Field) Nullable() bool {
	return field.cfg.Nullable
}
