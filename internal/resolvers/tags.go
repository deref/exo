package resolvers

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/graph-gophers/graphql-go/decode"
)

type Tags map[string]string

var _ decode.Unmarshaler = &Tags{}
var _ json.Marshaler = Tags{}
var _ sql.Scanner = &Tags{}
var _ driver.Valuer = Tags{}

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
