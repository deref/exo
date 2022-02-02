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

type JSONObject map[string]interface{}

var _ decode.Unmarshaler = &JSONObject{}
var _ json.Marshaler = JSONObject{}
var _ sql.Scanner = &JSONObject{}
var _ driver.Valuer = JSONObject{}

func (_ JSONObject) ImplementsGraphQLType(name string) bool {
	return name == "JSONObject"
}

func (obj *JSONObject) UnmarshalGraphQL(input interface{}) (err error) {
	m, ok := input.(map[string]interface{})
	if !ok {
		return errors.New("expected JSON object")
	}
	*obj = m
	return nil
}

func (obj *JSONObject) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", src)
	}
	return json.Unmarshal([]byte(s), (*map[string]interface{})(obj))
}

func (obj JSONObject) Value() (driver.Value, error) {
	return jsonutil.MarshalString((map[string]interface{})(obj))
}

func (obj JSONObject) MarshalJSON() ([]byte, error) {
	return json.Marshal((map[string]interface{})(obj))
}
