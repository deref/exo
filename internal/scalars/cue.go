package scalars

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/parser"
	"github.com/deref/exo/internal/util/cueutil"
)

type CueValue cue.Value

var _ Scalar = &CueValue{}

func (_ CueValue) ImplementsGraphQLType(name string) bool {
	return name == "CueValue"
}

func (cv *CueValue) UnmarshalGraphQL(input any) (err error) {
	return cv.unmarshal(input)
}

func (cv *CueValue) UnmarshalJSON(bs []byte) error {
	var v any
	if err := json.Unmarshal(bs, &v); err != nil {
		return err
	}
	return cv.unmarshal(v)
}

func (cv *CueValue) Scan(src any) error {
	return cv.unmarshal(src)
}

func (cv *CueValue) unmarshal(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}
	cc := cuecontext.New()
	f, err := parser.ParseFile("", s, parser.ParseComments)
	if err != nil {
		return err
	}
	underlying := cc.BuildFile(f)
	*cv = CueValue(underlying)
	return underlying.Err()
}

func (cv CueValue) Value() (driver.Value, error) {
	return cv.String(), nil
}

func (cv CueValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(cv.String())
}

func (cv CueValue) String() string {
	return cueutil.MustValueToString(cue.Value(cv))
}

func (cv CueValue) Bytes() []byte {
	return cueutil.MustValueToBytes(cue.Value(cv))
}

func EncodeCueValue(v any) CueValue {
	cc := cuecontext.New()
	return CueValue(cc.Encode(v))
}
