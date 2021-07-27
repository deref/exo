package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

type Serde interface {
	Serialize(*Tuple) ([]byte, error)
	Deserialize([]byte) (*Tuple, error)
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

	// Map of the the variable-length element index to the location where the offset
	// needs to be written.
	varlenIndexes := make(map[int]int)

	var buf bytes.Buffer
	for idx, elem := range tup.elements {
		elemSchema := s.schema.Elements[idx]

		switch elemSchema.Type {
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
}

func (s SchematizedRowSerde) Deserialize(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}
