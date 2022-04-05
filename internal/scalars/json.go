package scalars

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

// JSON scalar for values known to be map-like objects. Marshals to and from
// the database as an encoded JSON string, but marshals to and from GraphQL as
// plain JSON-values, avoiding double-encoding.
type JSONObject map[string]any

// JSON scalar with json.RawMessage style unmarshaling behavior.
type RawJSON []byte

var _ Scalar = &JSONObject{}

// Could reasonably be extended to full-Scalar support if needed.
var _ DatabaseScalar = &RawJSON{}

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

func (raw *RawJSON) Scan(src any) error {
	switch s := src.(type) {
	case nil:
		*raw = nil
		return nil
	case string:
		*raw = []byte(s)
		return nil
	default:
		return fmt.Errorf("expected string, got %T", src)
	}
}

func (obj JSONObject) Value() (driver.Value, error) {
	return jsonutil.MarshalString((map[string]any)(obj))
}

func (raw RawJSON) Value() (driver.Value, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return raw, nil
}

func (obj *JSONObject) UnmarshalJSON(bs []byte) (err error) {
	return json.Unmarshal(bs, (*map[string]any)(obj))
}

func (obj JSONObject) MarshalJSON() ([]byte, error) {
	return json.Marshal((map[string]any)(obj))
}
