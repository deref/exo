package storage

import (
	"encoding/json"
	"strings"

	"github.com/deref/exo/util/cmdutil"
)

type ElemType int

const (
	TypeInt32 ElemType = iota
	TypeUint32
	TypeInt64
	TypeUint64
	TypeBoolean
	TypeBytes
	TypeUnicode
	TypeUnknown
)

func (t ElemType) String() string {
	switch t {
	case TypeInt32:
		return "int32"

	case TypeUint32:
		return "uint32"

	case TypeInt64:
		return "int64"

	case TypeUint64:
		return "uint64"

	case TypeUnicode:
		return "unicode"

	case TypeBoolean:
		return "boolean"

	case TypeBytes:
		return "bytes"

	case TypeUnknown:
		return "unknown"

	default:
		return "<invalid>"
	}
}

func (t ElemType) DefaultValue() interface{} {
	switch t {
	case TypeInt32:
		return int32(0)

	case TypeUint32:
		return uint32(0)

	case TypeInt64:
		return int64(0)

	case TypeUint64:
		return uint64(0)

	case TypeUnicode:
		return ""

	case TypeBoolean:
		return false

	case TypeBytes:
		return []byte{}

	case TypeUnknown:
		return nil

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

func MustSerializeSchema(s *Schema) []byte {
	data, err := json.Marshal(s)
	if err != nil {
		cmdutil.Fatalf("cannot marshal schema: %w", err)
	}
	return data
}

func MustDeserializeSchema(data []byte) *Schema {
	s := &Schema{}
	if err := json.Unmarshal(data, s); err != nil {
		cmdutil.Fatalf("cannot unmarshal schema: %w", err)
	}
	return s
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

func (s *Schema) Slice(start, end int) *Schema {
	newElements := make([]ElementDescriptor, end-start)
	copy(newElements, s.Elements[start:end])
	return &Schema{Elements: newElements}
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
