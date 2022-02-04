package scalars

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

type ULID [16]byte

var _ Scalar = &ULID{}

func (_ ULID) ImplementsGraphQLType(name string) bool {
	return name == "ULID"
}

func (u *ULID) UnmarshalGraphQL(input interface{}) error {
	return u.unmarshal(input)
}

func (u *ULID) Scan(src interface{}) error {
	return u.unmarshal(src)
}

func (u *ULID) unmarshal(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}
	parsed, err := ulid.ParseStrict(s)
	if err != nil {
		return err
	}
	*u = ([16]byte)(parsed)
	return nil
}

func (u ULID) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u ULID) UnmarshalJSON(bs []byte) (err error) {
	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return err
	}
	return u.unmarshal(s)
}

func (u ULID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u ULID) String() string {
	return ulid.ULID(u).String()
}

func (u ULID) Timestamp() Instant {
	ms := int64(ulid.ULID([16]byte(u)).Time())
	return Instant{time.Unix(0, ms*int64(time.Millisecond))}
}

// Greater than any valid ULID.
var InfiniteULID = ULID{
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
}

func ULIDMin(a, b ULID) ULID {
	if bytes.Compare(a[:], b[:]) < 0 {
		return a
	}
	return b
}

func ULIDMax(a, b ULID) ULID {
	if bytes.Compare(a[:], b[:]) < 0 {
		return b
	}
	return a
}

func IncrementULID(id ULID) ULID {
	bs := id[:]
	res := id
	for i := len(bs) - 1; i >= 0; i-- {
		b := bs[i]
		if b == 255 {
			res[i] = 0
		} else {
			res[i] = b + 1
			return ULID(res)
		}
	}
	panic("overflow")
}
