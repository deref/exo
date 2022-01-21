package peer

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/graph-gophers/graphql-go"
	"github.com/jmoiron/sqlx"
)

type Peer struct {
	db     *sqlx.DB
	schema *graphql.Schema
}

var _ api.Service = (*Peer)(nil)

func NewPeer(ctx context.Context, varDir string) (*Peer, error) {
	// XXX reconsile SQL DB opening with exod/main.go
	dbPath := filepath.Join(varDir, "exo.sqlite3")
	txMode := "exclusive"
	connStr := dbPath + "?_txlock=" + txMode
	db, err := sqlx.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db: %w", err)
	}
	r := &resolvers.RootResolver{
		DB:     db,
		Logger: logging.CurrentLogger(ctx),
	}
	// XXX migration probably shouldn't happen here.
	if err := r.Migrate(ctx); err != nil {
		cmdutil.Fatalf("migrating db: %w", err)
	}
	return &Peer{
		db:     db,
		schema: resolvers.NewSchema(r),
	}, nil
}

func (p *Peer) Do(ctx context.Context, doc string, vars map[string]interface{}, res interface{}) error {
	vars, _ = jsonutil.MustSimplify(vars).(map[string]interface{})

	resp := p.schema.Exec(ctx, doc, "", vars)
	if len(resp.Errors) > 0 {
		return api.QueryErrorSet(resp.Errors)
	}
	if err := json.Unmarshal(resp.Data, res); err != nil {
		return err
	}
	return nil
}

func (p *Peer) Shutdown(ctx context.Context) error {
	if err := p.db.Close(); err != nil {
		return fmt.Errorf("closing sqlite db: %w", err)
	}
	return nil
}
