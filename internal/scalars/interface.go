package scalars

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/graph-gophers/graphql-go/decode"
)

type Scalar interface {
	decode.Unmarshaler

	json.Marshaler
	json.Unmarshaler

	sql.Scanner
	driver.Valuer
}
