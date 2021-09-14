package atom

import (
	"encoding/json"
	"fmt"
)

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(bs []byte, v interface{}) error
}

type jsonCodec struct{}

var CodecJSON = &jsonCodec{}

func (codec *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	bs, err := json.MarshalIndent(v, "", "  ")
	bs = append(bs, '\n')
	return bs, err
}

func (codec *jsonCodec) Unmarshal(bs []byte, v interface{}) error {
	return json.Unmarshal(bs, v)
}

type stringCodec struct{}

var CodecString = &stringCodec{}

func (codec *stringCodec) Marshal(v interface{}) ([]byte, error) {
	switch s := v.(type) {
	case string:
		return []byte(s), nil
	case *string:
		return []byte(*s), nil
	default:
		return nil, fmt.Errorf("string codec can only marshal string input but got %T", v)
	}
}

func (codec *stringCodec) Unmarshal(bs []byte, v interface{}) error {
	asStringPtr, ok := v.(*string)
	if !ok {
		return fmt.Errorf("string codec can only unmarshal to string pointer but got %T", v)
	}
	*asStringPtr = string(bs)
	return nil
}
