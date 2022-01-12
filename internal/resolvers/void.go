package resolvers

import (
	"encoding/json"
	"errors"

	"github.com/graph-gophers/graphql-go/decode"
)

type Void struct{}

var _ decode.Unmarshaler = Void{}
var _ json.Marshaler = Void{}

func (_ Void) ImplementsGraphQLType(name string) bool {
	return name == "Void"
}

func (_ Void) UnmarshalGraphQL(input interface{}) error {
	if input != nil {
		return errors.New("invalid Void value")
	}
	return nil
}

func (_ Void) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
