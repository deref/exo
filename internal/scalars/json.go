package scalars

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

type JSONObject map[string]any

var _ Scalar = &JSONObject{}

func (_ JSONObject) ImplementsGraphQLType(name string) bool {
	return name == "JSONObject"
}

func (obj *JSONObject) UnmarshalGraphQL(input any) (err error) {
	m, ok := input.(map[string]any)
	if !ok {
		return errors.New("expected JSON object")
	}
	*obj = m
	return nil
}

func (obj *JSONObject) Scan(src any) error {
	switch s := src.(type) {
	case nil:
		*obj = nil
		return nil
	case string:
		return json.Unmarshal([]byte(s), (*map[string]any)(obj))
	default:
		return fmt.Errorf("expected string, got %T", src)
	}
}

func (obj JSONObject) Value() (driver.Value, error) {
	return jsonutil.MarshalString((map[string]any)(obj))
}

func (obj *JSONObject) UnmarshalJSON(bs []byte) (err error) {
	return json.Unmarshal(bs, (*map[string]any)(obj))
}

func (obj JSONObject) MarshalJSON() ([]byte, error) {
	return json.Marshal((map[string]any)(obj))
}
