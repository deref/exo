package resolvers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
var _ driver.Valuer = Instant{}

func Now(ctx context.Context) Instant {
	return Instant{chrono.Now(ctx)}
}

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
		return fmt.Errorf("expected string, got %T", v)
	}
	inst.t, err = chrono.ParseIsoNano(s)
	return
}

func (inst Instant) Value() (driver.Value, error) {
	return inst.String(), nil
}

func (inst Instant) MarshalJSON() ([]byte, error) {
	return json.Marshal(inst.String())
}

func (inst Instant) String() string {
	return chrono.IsoNano(inst.t)
}
