package resolvers

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/logging"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jmoiron/sqlx"
)

//go:embed schema.gql
var schema string

func NewSchema(r *RootResolver) *graphql.Schema {
	return graphql.MustParseSchema(schema, r,
		graphql.UseFieldResolvers(),
		graphql.Logger(NewGraphqlLogger(r.SystemLog)),
	)
}

type RootResolver struct {
	SystemLog   logging.Logger
	VarDir      string
	GUIEndpoint string

	ulidgen *gensym.ULIDGenerator
	db      *sqlx.DB
}

func (r *RootResolver) Init(ctx context.Context) error {
	r.ulidgen = gensym.NewULIDGenerator(ctx)

	dbPath := filepath.Join(r.VarDir, "exo.sqlite3")
	var err error
	r.db, err = OpenDB(ctx, dbPath)
	if err != nil {
		return fmt.Errorf("opening sqlite db: %w", err)
	}

	if err := r.Migrate(ctx); err != nil {
		return fmt.Errorf("migrating db: %w", err)
	}

	return nil
}

func (r *RootResolver) Shutdown(ctx context.Context) error {
	if err := r.db.Close(); err != nil {
		return fmt.Errorf("closing sqlite db: %w", err)
	}

	return nil
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
