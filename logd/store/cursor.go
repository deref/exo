package store

import (
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/deref/exo/util/bitutil"
)

var encoding32 = base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz").WithPadding(base32.NoPadding)

type Cursor struct {
	ID        []byte
	Direction direction
}

func (c Cursor) IsValid() bool {
	return c.ID != nil
}

const (
	cursorFlagForward byte = 0
)

func ParseCursor(in string) (*Cursor, error) {
	if len(in) == 0 {
		return nil, nil
	}
	if len(in) == 1 {
		return nil, errors.New("invalid cursor")
	}

	buf, err := encoding32.DecodeString(in)
	if err != nil {
		return nil, fmt.Errorf("decoding cursor: %w", err)
	}
	tag := buf[0]
	data := buf[1:]

	isForward := bitutil.FlagSetInByte(tag, cursorFlagForward)

	// TODO: Validate that data is valid for ID.
	c := &Cursor{
		ID: data,
	}

	if isForward {
		c.Direction = DirectionForward
	} else {
		c.Direction = DirectionReverse
	}

	return c, nil
}

func (c Cursor) Serialize() string {
	out := make([]byte, len(c.ID)+1)
	copy(out[1:], c.ID)

	var tag byte
	if c.Direction == DirectionForward {
		tag |= (1 << cursorFlagForward)
	}
	out[0] = tag

	return encoding32.EncodeToString(out)
}

type direction int8

const (
	DirectionForward direction = 1
	DirectionReverse direction = -1
)
