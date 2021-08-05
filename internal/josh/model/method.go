package model

type Method struct {
	iface   *Interface
	name    string
	doc     *string
	inputs  []*Field
	outputs []*Field
}

func (method *Method) Name() string {
	return method.name
}

func (method *Method) SetDoc(value *string) {
	method.doc = value
}

func (method *Method) Doc() *string {
	return method.doc
}

func (method *Method) Inputs() []*Field {
	inputs := make([]*Field, len(method.inputs))
	copy(inputs, method.inputs)
	return inputs
}

func (method *Method) Outputs() []*Field {
	outputs := make([]*Field, len(method.outputs))
	copy(outputs, method.outputs)
	return outputs
}

func (method *Method) AddInput(cfg FieldConfig) *Field {
	field := &Field{cfg: cfg}
	method.inputs = append(method.inputs, field)
	return field
}

func (method *Method) AddOutput(cfg FieldConfig) *Field {
	field := &Field{cfg: cfg}
	method.outputs = append(method.outputs, field)
	return field
}
