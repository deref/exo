package scalars

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	graphql "github.com/graph-gophers/graphql-go/decode"
)

type Scalar interface {
	graphql.Unmarshaler

	json.Marshaler
	json.Unmarshaler

	sql.Scanner
	driver.Valuer
}
