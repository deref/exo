package storage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

var ErrOutOfRange = errors.New("index out of range")

type Tuple struct {
	elements   []interface{}
	schema     *Schema
	serialized []byte
}

func NewTuple(elements ...interface{}) *Tuple {
	t := &Tuple{schema: NewSchema()}

	for _, elem := range elements {
		switch elem.(type) {
		case string:
			t.schema.AppendUnnamed(TypeUnicode)
		case int64:
			t.schema.AppendUnnamed(TypeInt64)
		case uint64:
			t.schema.AppendUnnamed(TypeUint64)
		default:
			t.schema.AppendUnnamed(TypeUnknown)
		}
		t.elements = append(t.elements, elem)
	}

	return t
}

func NewTupleWithSchema(schema *Schema) *Tuple {
	return &Tuple{
		schema:   schema,
		elements: make([]interface{}, len(schema.Elements)),
	}
}

func (t *Tuple) Schema() *Schema {
	return t.schema
}

func (t *Tuple) Concat(other *Tuple) *Tuple {
	return &Tuple{
		elements: append(append([]interface{}{}, t.elements...), other.elements...),
		schema:   t.schema.Concat(other.schema),
	}
}

func (t *Tuple) Without(idx int) *Tuple {
	return &Tuple{
		elements: append(append([]interface{}{}, t.elements[:idx]...), t.elements[idx+1:]...),
		schema:   t.schema.Without(idx),
	}
}

func (t *Tuple) Size() int {
	return len(t.elements)
}

func (t *Tuple) GetUntyped(i int) (interface{}, error) {
	if i > len(t.elements)-1 {
		return "", fmt.Errorf("Index out of range: %d", i)
	}

	return t.elements[i], nil
}

// TODO: Codegen these.

func (t *Tuple) GetUnicode(i int) (string, error) {
	if i > len(t.elements)-1 {
		return "", ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeUnicode {
		return "", fmt.Errorf("Expected string at %d but got %s", i, elemType)
	}

	return t.elements[i].(string), nil
}

func (t *Tuple) SetUnicode(i int, val string) error {
	if i > len(t.elements)-1 {
		return ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeUnicode {
		return fmt.Errorf("Expected string at %d but got %s", i, elemType)
	}

	t.elements[i] = val
	return nil
}

func (t *Tuple) GetInt64(i int) (int64, error) {
	if i > len(t.elements)-1 {
		return 0, ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeInt64 {
		return 0, fmt.Errorf("Expected int64 at %d but got %s", i, elemType)
	}

	return t.elements[i].(int64), nil
}

func (t *Tuple) SetInt64(i int, val int64) error {
	if i > len(t.elements)-1 {
		return ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeInt64 {
		return fmt.Errorf("Expected int64 at %d but got %s", i, elemType)
	}

	t.elements[i] = val
	return nil
}

func (t *Tuple) GetUint64(i int) (uint64, error) {
	if i > len(t.elements)-1 {
		return 0, ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeUint64 {
		return 0, fmt.Errorf("Expected uint64 at %d but got %s", i, elemType)
	}

	return t.elements[i].(uint64), nil
}

func (t *Tuple) SetUint64(i int, val uint64) error {
	if i > len(t.elements)-1 {
		return ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeUint64 {
		return fmt.Errorf("Expected uint64 at %d but got %s", i, elemType)
	}

	t.elements[i] = val
	return nil
}

func (t *Tuple) Serialize() []byte {
	if t.serialized == nil {
		t.doSerialize()
	}
	return t.serialized
}

func (t *Tuple) doSerialize() {
	buf := bytes.NewBuffer(make([]byte, 0, t.sizeHint()))

	for i, elementDescriptor := range t.schema.Elements {
		typ := elementDescriptor.Type
		buf.WriteByte(byte(typ))
		elem := t.elements[i]
		switch typ {
		case TypeUnicode:
			buf.Write(escapeNulls([]byte(elem.(string))))
			buf.WriteByte(0)

		case TypeInt64:
			buf := bytes.NewBuffer(make([]byte, 0, 8))
			_ = binary.Write(buf, binary.BigEndian, elem.(int64))
			buf.Write(buf.Bytes())

		case TypeUint64:
			i := make([]byte, 8)
			binary.BigEndian.PutUint64(i, elem.(uint64))
			buf.Write(i)

		default:
			panic(fmt.Errorf("no serializer defined for %s@%d", elementDescriptor, i))
		}
	}

	t.serialized = buf.Bytes()
}

// Deserialize takes a byte slice and attempts to deserialize it into a Tuple. If the
// slice does not contain a valid serialzed Tuple, an error will be returned.
// NOTE: the tuple will retain a reference to `buf`, so the underlying data must not
// be mutated after Deserialization.
func Deserialize(buf []byte) (*Tuple, error) {
	t := &Tuple{schema: NewSchema()}

	buflen := len(buf)
	for byteIdx := 0; byteIdx < buflen; byteIdx++ {
		typeTag := ElemType(buf[byteIdx])
		t.schema.AppendUnnamed(typeTag)
		// Skip past type tag
		byteIdx++

		switch typeTag {
		case TypeUnicode:
			var s []byte
			for maybeEscapedBytes := 0; ; maybeEscapedBytes++ {
				ch := buf[byteIdx+maybeEscapedBytes]
				if ch == 0x00 {
					if byteIdx+maybeEscapedBytes+1 < buflen && buf[byteIdx+maybeEscapedBytes+1] == 0xff {
						// The null byte was escaped, so we allow the null byte to be emitted
						s = append(s, ch)
						maybeEscapedBytes++
						continue
					}
					// Advance iterator by characters consumed.
					byteIdx += maybeEscapedBytes
					t.elements = append(t.elements, string(s))
					break
				}
				s = append(s, ch)
			}

		case TypeInt64:
			var x int64
			if err := binary.Read(bytes.NewReader(buf[byteIdx:byteIdx+8]), binary.BigEndian, &x); err != nil {
				return nil, fmt.Errorf("error decoding value as int64")
			}
			byteIdx += 7

		case TypeUint64:
			t.elements = append(t.elements, binary.BigEndian.Uint64(buf[byteIdx:byteIdx+8]))
			byteIdx += 7

		default:
			return nil, fmt.Errorf("Invalid type tag at %d: %x", byteIdx, typeTag)
		}
	}

	t.serialized = buf

	return t, nil
}

func (t *Tuple) String() string {
	var sb strings.Builder
	sb.WriteByte('(')
	for i, elem := range t.elements {
		typ := t.schema.MustGet(i).Type
		switch typ {
		case TypeUnicode:
			sb.WriteString(elem.(string))
		case TypeInt64:
			sb.WriteString(fmt.Sprintf("%d", elem))
		case TypeUint64:
			sb.WriteString(fmt.Sprintf("%d", elem))
		default:
			return "<noprint>"
		}
		sb.WriteByte(':')
		sb.WriteString(typ.String())

		if i < len(t.elements)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteByte(')')

	return sb.String()
}

func (t *Tuple) sizeHint() int {
	size := 0
	for i, elementDescriptor := range t.schema.Elements {
		size++ // type tag
		switch elementDescriptor.Type {
		case TypeUnicode:
			size += len(t.elements[i].(string))
			size++ // delimiter
		case TypeUint64:
			size += 8
		case TypeInt64:
			size += 8
		}
	}

	return size
}

// escapeNulls retuns a version of a byte string where all instances of a null
// byte (0x00) are followed by a 0xff byte.
func escapeNulls(in []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, len(in)))

	for _, byt := range in {
		buf.WriteByte(byt)
		if byt == 0x00 {
			buf.WriteByte(0xff)
		}
	}

	return buf.Bytes()
}
