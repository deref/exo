package resolvers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/graph-gophers/graphql-go/decode"
	"github.com/oklog/ulid/v2"
)

type ULID [16]byte

var _ decode.Unmarshaler = &ULID{}
var _ json.Marshaler = ULID{}
var _ sql.Scanner = &ULID{}
var _ driver.Valuer = ULID{}

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
		return errors.New("expected string")
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

func (u ULID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u ULID) String() string {
	return ulid.ULID(u).String()
}

func (r *MutationResolver) mustNextULID(ctx context.Context) ULID {
	res, err := r.ULIDGenerator.NextID(ctx)
	if err != nil {
		panic(err)
	}
	return ULID(res)
}
