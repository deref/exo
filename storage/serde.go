package storage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Serde interface {
	Serialize(*Tuple) ([]byte, error)
	Deserialize([]byte) (*Tuple, error)
}

func NewSchematizedRowSerde(schema *Schema) *SchematizedRowSerde {
	return &SchematizedRowSerde{
		schema: schema,
	}
}

type SchematizedRowSerde struct {
	schema *Schema
}

func (s SchematizedRowSerde) Serialize(tup *Tuple) ([]byte, error) {
	if s.schema != tup.Schema() {
		return nil, errors.New("row schema does not match serializer")
	}

	// NOTE [VARIABLE LENGTH DATA]:
	// For variable-length data, we write the offset and length of the data in
	// the main body of the tuple so that we can jump to any fixed-length element
	// in constant time. Reading a variable-length field involves jumping to its
	// element location in the main tuple, reading the offset and length of the
	// variable data, then reading the variable data at the location indicated.
	// In order to do this, we need to keep track of where in the main storage
	// we need to write the offset, then after appending all of the variable-length
	// data, we go back and write the offsets afterwards.

	type placeholderLocation struct {
		elemIdx   int
		offsetLoc int
	}
	placeholderLocations := make([]placeholderLocation, 0, len(tup.elements))

	var buf bytes.Buffer
	for idx, elem := range tup.elements {
		elemSchema := s.schema.Elements[idx]

		switch elemSchema.Type {
		case TypeUnicode:
			placeholderLocations = append(placeholderLocations, placeholderLocation{
				elemIdx:   idx,
				offsetLoc: buf.Len(),
			})
			// Write placeholder to be replaced with 32-bit offset.
			buf.Write([]byte{0, 0, 0, 0})

			// Write 32-bit length.
			length := make([]byte, 4)
			strLen := len(elem.(string))
			binary.BigEndian.PutUint32(length, uint32(strLen))
			buf.Write(length)

		case TypeInt64:
			n := make([]byte, 8)
			binary.BigEndian.PutUint64(n, uint64(elem.(int64)))
			buf.Write(n)

		case TypeUint64:
			i := make([]byte, 8)
			binary.BigEndian.PutUint64(i, elem.(uint64))
			buf.Write(i)

		default:
			panic(fmt.Errorf("no serializer defined for %s@%d", elem, idx))
		}
	}

	type offsetReference struct {
		refAt  int
		dataAt int
	}
	offsets := make([]offsetReference, 0, len(placeholderLocations))
	for _, placeholderLocation := range placeholderLocations {
		elemIdx := placeholderLocation.elemIdx
		offsetLoc := placeholderLocation.offsetLoc
		// TODO: handle bytes as well.
		var varlenData []byte
		elem := tup.elements[elemIdx]
		elemSchema := s.schema.Elements[elemIdx]
		switch elemSchema.Type {
		case TypeUnicode:
			varlenData = []byte(elem.(string))
		default:
			panic("unhandled variable-length type")
		}
		dataLoc := buf.Len()
		offsets = append(offsets, struct {
			refAt  int
			dataAt int
		}{
			refAt:  offsetLoc,
			dataAt: dataLoc,
		})
		buf.Write(varlenData)
	}

	out := buf.Bytes()

	// Patch the data by writing the offsets back to the buffer.
	for _, offset := range offsets {
		start := offset.refAt
		binary.BigEndian.PutUint32(out[start:start+4], uint32(offset.dataAt))
	}

	return out, nil
}

func (s SchematizedRowSerde) Deserialize(buf []byte) (*Tuple, error) {
	t := &Tuple{
		schema:   s.schema,
		elements: make([]interface{}, 0, len(s.schema.Elements)),
	}

	var pos int
	for _, elemSchema := range s.schema.Elements {
		switch elemSchema.Type {
		case TypeInt64:
			n := binary.BigEndian.Uint64(buf[pos : pos+8])
			t.elements = append(t.elements, int64(n))
			pos += 8

		case TypeUint64:
			n := binary.BigEndian.Uint64(buf[pos : pos+8])
			t.elements = append(t.elements, n)
			pos += 8

		case TypeUnicode:
			dataOffset := binary.BigEndian.Uint32(buf[pos : pos+4])
			pos += 4

			dataLen := binary.BigEndian.Uint32(buf[pos : pos+4])
			pos += 4

			str := string(buf[dataOffset : dataOffset+dataLen])
			t.elements = append(t.elements, str)

		default:
			panic("Cannot handle type for deserialization")
		}
	}

	return t, nil
}
