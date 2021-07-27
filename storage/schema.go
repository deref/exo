package storage

import "strings"

type ElemType int

const (
	TypeUnicode ElemType = iota
	TypeInt64
	TypeUint64
	TypeUnknown
)

func (t ElemType) String() string {
	switch t {
	case TypeUnicode:
		return "string"
	case TypeInt64:
		return "int64"
	case TypeUint64:
		return "uint64"
	case TypeUnknown:
		return "unknown"
	default:
		return "<invalid>"
	}
}

type Schema struct {
	Elements []ElementDescriptor
}

func NewSchema(elements ...ElementDescriptor) *Schema {
	return &Schema{Elements: elements}
}

func (s *Schema) Append(typ ElemType, name string) {
	s.Elements = append(s.Elements, ElementDescriptor{
		Type: typ,
		Name: name,
	})
}

func (s *Schema) AppendUnnamed(typ ElemType) {
	s.Elements = append(s.Elements, ElementDescriptor{
		Type: typ,
	})
}

func (s *Schema) Concat(other *Schema) *Schema {
	elements := make([]ElementDescriptor, len(s.Elements)+len(other.Elements))
	var i int
	for _, elem := range s.Elements {
		elements[i] = elem
		i++
	}
	for _, elem := range other.Elements {
		elements[i] = elem
		i++
	}

	return &Schema{Elements: elements}
}

func (s *Schema) Without(idx int) *Schema {
	elements := append(append([]ElementDescriptor{}, s.Elements[:idx]...), s.Elements[idx+1:]...)

	return &Schema{Elements: elements}
}

func (s *Schema) Get(idx int) (ElementDescriptor, bool) {
	if idx >= len(s.Elements) {
		return ElementDescriptor{}, false
	}

	return s.Elements[idx], true
}

func (s *Schema) MustGet(idx int) ElementDescriptor {
	return s.Elements[idx]
}

func (s *Schema) GetNamed(name string) (ElementDescriptor, bool) {
	for _, element := range s.Elements {
		if element.Name == name {
			return element, true
		}
	}

	return ElementDescriptor{}, false
}

type ElementDescriptor struct {
	Type ElemType
	Name string
}

func (e ElementDescriptor) String() string {
	var sb strings.Builder
	if e.Name == "" {
		sb.WriteByte('_')
	} else {
		sb.WriteString(e.Name)
	}
	sb.WriteByte(':')
	sb.WriteString(e.Type.String())

	return sb.String()
}
