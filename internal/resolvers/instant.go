package resolvers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/deref/exo/internal/chrono"
	"github.com/graph-gophers/graphql-go/decode"
)

type Instant struct {
	t time.Time
}

var _ decode.Unmarshaler = &Instant{}
var _ json.Marshaler = Instant{}
var _ sql.Scanner = &Instant{}

func (_ Instant) ImplementsGraphQLType(name string) bool {
	return name == "Instant"
}

func (inst *Instant) UnmarshalGraphQL(input interface{}) (err error) {
	return inst.unmarshal(input)
}

func (inst *Instant) Scan(src interface{}) error {
	return inst.unmarshal(src)
}

func (inst *Instant) unmarshal(v interface{}) (err error) {
	s, ok := v.(string)
	if !ok {
		return errors.New("expected string")
	}
	inst.t, err = chrono.ParseIsoNano(s)
	return
}

func (inst Instant) MarshalJSON() ([]byte, error) {
	return []byte(inst.String()), nil
}

func (inst Instant) String() string {
	return chrono.IsoNano(inst.t)
}
