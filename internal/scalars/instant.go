package scalars

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deref/exo/internal/chrono"
)

type Instant struct {
	t time.Time
}

var _ Scalar = &Instant{}

func GoTimeToInstant(t time.Time) Instant {
	return Instant{t.UTC()}

}
func (inst Instant) GoTime() time.Time {
	return inst.t
}

func Now(ctx context.Context) Instant {
	return Instant{chrono.Now(ctx)}
}

func (_ Instant) ImplementsGraphQLType(name string) bool {
	return name == "Instant"
}

func (inst *Instant) UnmarshalGraphQL(input interface{}) (err error) {
	return inst.unmarshal(input)
}

func (inst *Instant) UnmarshalJSON(bs []byte) (err error) {
	var v interface{}
	if err := json.Unmarshal(bs, &v); err != nil {
		return err
	}
	return inst.unmarshal(v)
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

func (inst Instant) UnixMilli() int64 {
	return inst.t.UnixMilli()
}

func (inst Instant) Before(other Instant) bool {
	return inst.GoTime().Before(other.GoTime())
}

func (inst Instant) After(other Instant) bool {
	return inst.GoTime().After(other.GoTime())
}

func (inst Instant) Equal(other Instant) bool {
	return inst.GoTime().Equal(other.GoTime())
}

func (inst Instant) Sub(other Instant) time.Duration {
	return inst.GoTime().Sub(other.GoTime())
}
