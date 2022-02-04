package scalars

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

type Tags map[string]string

var _ Scalar = &Tags{}

func (_ Tags) ImplementsGraphQLType(name string) bool {
	return name == "Tags"
}

func (t *Tags) UnmarshalGraphQL(input interface{}) (err error) {
	m, ok := input.(map[string]interface{})
	if !ok {
		return errors.New("expected JSON object")
	}
	*t = make(map[string]string, len(m))
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("expected key %q to have string value", k)
		}
		(*t)[k] = s
	}
	return nil
}

func (t *Tags) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", src)
	}
	return json.Unmarshal([]byte(s), (*map[string]string)(t))
}

func (t Tags) Value() (driver.Value, error) {
	return jsonutil.MarshalString((map[string]string)(t))
}

func (t Tags) MarshalJSON() ([]byte, error) {
	return json.Marshal((map[string]string)(t))
}

func (t *Tags) UnmarshalJSON(bs []byte) (err error) {
	return json.Unmarshal(bs, (*map[string]string)(t))
}
