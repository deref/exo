package atom

import (
	"encoding/json"
	"fmt"
)

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(bs []byte, v any) error
}

type jsonCodec struct{}

var CodecJSON = &jsonCodec{}

func (codec *jsonCodec) Marshal(v any) ([]byte, error) {
	bs, err := json.MarshalIndent(v, "", "  ")
	bs = append(bs, '\n')
	return bs, err
}

func (codec *jsonCodec) Unmarshal(bs []byte, v any) error {
	return json.Unmarshal(bs, v)
}

type stringCodec struct{}

var CodecString = &stringCodec{}

func (codec *stringCodec) Marshal(v any) ([]byte, error) {
	switch s := v.(type) {
	case string:
		return []byte(s), nil
	case *string:
		return []byte(*s), nil
	default:
		return nil, fmt.Errorf("string codec can only marshal string input but got %T", v)
	}
}

func (codec *stringCodec) Unmarshal(bs []byte, v any) error {
	asStringPtr, ok := v.(*string)
	if !ok {
		return fmt.Errorf("string codec can only unmarshal to string pointer but got %T", v)
	}
	*asStringPtr = string(bs)
	return nil
}
