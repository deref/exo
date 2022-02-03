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
	SystemLog     logging.Logger
	ULIDGenerator *gensym.ULIDGenerator
	Routes        *RoutesResolver
}

// While queries, mutations, and subscriptions are accessed in disjoint query
// paths, this GraphQL library assumes that their names will not conflict and
// therefore all resolvers go on the same struct. We use the following aliases
// for clarity. See <https://github.com/graph-gophers/graphql-go/pull/182> for
// more details. Note that even if this were not required, it's still
// convenient to be able to access query methods from mutations and
// subscriptions.
type QueryResolver = RootResolver
type MutationResolver = RootResolver
type SubscriptionResolver = RootResolver
