package store

import (
	"encoding/base64"
)

var cursorEnc = base64.RawURLEncoding

var NilCursor = &Cursor{
	ID: make([]byte, 16),
}

type Cursor struct {
	ID []byte
}

func ParseCursor(in string) (Cursor, error) {
	data, err := cursorEnc.DecodeString(in)
	if err != nil {
		return Cursor{}, err
	}

	return Cursor{
		ID: data,
	}, nil
}

func (c *Cursor) Serialize() string {
	return cursorEnc.EncodeToString(c.ID)
}
