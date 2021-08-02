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
		case int32:
			t.schema.AppendUnnamed(TypeInt32)
		case uint32:
			t.schema.AppendUnnamed(TypeUint32)
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

func (t *Tuple) GetInt32(i int) (int32, error) {
	if i > len(t.elements)-1 {
		return 0, ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeInt32 {
		return 0, fmt.Errorf("Expected int32 at %d but got %s", i, elemType)
	}

	return t.elements[i].(int32), nil
}

func (t *Tuple) SetInt32(i int, val int32) error {
	if i > len(t.elements)-1 {
		return ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeInt32 {
		return fmt.Errorf("Expected int32 at %d but got %s", i, elemType)
	}

	t.elements[i] = val
	return nil
}

func (t *Tuple) GetUint32(i int) (uint32, error) {
	if i > len(t.elements)-1 {
		return 0, ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeUint32 {
		return 0, fmt.Errorf("Expected uint32 at %d but got %s", i, elemType)
	}

	return t.elements[i].(uint32), nil
}

func (t *Tuple) SetUint32(i int, val uint32) error {
	if i > len(t.elements)-1 {
		return ErrOutOfRange
	}

	elemType := t.schema.MustGet(i).Type
	if elemType != TypeUint32 {
		return fmt.Errorf("Expected uint32 at %d but got %s", i, elemType)
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
		case TypeInt32:
			_ = binary.Write(buf, binary.BigEndian, elem.(int32))

		case TypeUint32:
			_ = binary.Write(buf, binary.BigEndian, elem.(uint32))

		case TypeInt64:
			_ = binary.Write(buf, binary.BigEndian, elem.(int64))

		case TypeUint64:
			_ = binary.Write(buf, binary.BigEndian, elem.(uint64))

		case TypeBoolean:
			if elem.(bool) {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}

		case TypeBytes:
			buf.Write(escapeNulls(elem.([]byte)))
			buf.WriteByte(0)

		case TypeUnicode:
			buf.Write(escapeNulls([]byte(elem.(string))))
			buf.WriteByte(0)

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
	byteIdx := 0
	for byteIdx < buflen {
		typeTag := ElemType(buf[byteIdx])
		t.schema.AppendUnnamed(typeTag)
		// Skip past type tag
		byteIdx++

		switch typeTag {
		case TypeInt32:
			asUint32 := binary.BigEndian.Uint32(buf[byteIdx : byteIdx+4])
			t.elements = append(t.elements, int32(asUint32))
			byteIdx += 4

		case TypeUint32:
			t.elements = append(t.elements, binary.BigEndian.Uint32(buf[byteIdx:byteIdx+4]))
			byteIdx += 4

		case TypeInt64:
			asUint64 := binary.BigEndian.Uint64(buf[byteIdx : byteIdx+8])
			t.elements = append(t.elements, int64(asUint64))
			byteIdx += 8

		case TypeUint64:
			t.elements = append(t.elements, binary.BigEndian.Uint64(buf[byteIdx:byteIdx+8]))
			byteIdx += 8

		case TypeBoolean:
			t.elements = append(t.elements, buf[byteIdx] == 1)
			byteIdx += 1

		case TypeBytes, TypeUnicode:
			var bs []byte
			for maybeEscapedBytes := 0; ; maybeEscapedBytes++ {
				ch := buf[byteIdx+maybeEscapedBytes]
				if ch == 0x00 {
					if byteIdx+maybeEscapedBytes+1 < buflen && buf[byteIdx+maybeEscapedBytes+1] == 0xff {
						// The null byte was escaped, so we allow the null byte to be emitted
						bs = append(bs, ch)
						maybeEscapedBytes++
						continue
					}

					// Advance past characters consumed.
					byteIdx += maybeEscapedBytes
					// Advance past null terminator.
					byteIdx += 1

					if typeTag == TypeBytes {
						t.elements = append(t.elements, bs)
					} else {
						t.elements = append(t.elements, string(bs))
					}
					break
				}
				bs = append(bs, ch)
			}

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

		case TypeInt32, TypeUint32,
			TypeInt64, TypeUint64:

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
