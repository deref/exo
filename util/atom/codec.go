package atom

import "encoding/json"

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
