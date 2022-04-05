package scalars

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	graphql "github.com/graph-gophers/graphql-go/decode"
)

type Scalar interface {
	GraphQLScalar
	DatabaseScalar
}

type GraphQLScalar interface {
	graphql.Unmarshaler

	json.Marshaler
	json.Unmarshaler
}

type DatabaseScalar interface {
	sql.Scanner
	driver.Valuer
}
