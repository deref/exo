package resolvers

import (
	_ "embed"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/jmoiron/sqlx"
)

//go:embed schema.gql
var schema string

func NewHandler(resolver *RootResolver) *relay.Handler {
	return &relay.Handler{
		Schema: graphql.MustParseSchema(schema, resolver, graphql.UseFieldResolvers()),
	}
}

type RootResolver struct {
	DB *sqlx.DB
}

// While queries and mutations are accessed in disjoint query paths, this
// GraphQL library assumes that their names will not conflict and therefore all
// resolvers go on the same struct. We use the following aliases for clarity.
// See <https://github.com/graph-gophers/graphql-go/pull/182>.
type QueryResolver = RootResolver
type MutationResolver = RootResolver
