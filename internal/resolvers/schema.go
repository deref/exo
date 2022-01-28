package resolvers

import (
	_ "embed"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/logging"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/jmoiron/sqlx"
)

//go:embed schema.gql
var schema string

func NewSchema(r *RootResolver) *graphql.Schema {
	return graphql.MustParseSchema(schema, r, graphql.UseFieldResolvers())
}

// XXX move this to server package.
func NewHandler(r *RootResolver) *relay.Handler {
	return &relay.Handler{
		Schema: NewSchema(r),
	}
}

type RootResolver struct {
	DB            *sqlx.DB
	Logger        logging.Logger
	ULIDGenerator *gensym.ULIDGenerator
}

// While queries and mutations are accessed in disjoint query paths, this
// GraphQL library assumes that their names will not conflict and therefore all
// resolvers go on the same struct. We use the following aliases for clarity.
// See <https://github.com/graph-gophers/graphql-go/pull/182>.
type QueryResolver = RootResolver
type MutationResolver = RootResolver
