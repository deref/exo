package atom

import "encoding/json"

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(bs []byte, v interface{}) error
}

type jsonCodec struct{}

var CodecJSON = &jsonCodec{}

func (codec *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func (codec *jsonCodec) Unmarshal(bs []byte, v interface{}) error {
	return json.Unmarshal(bs, v)
}
