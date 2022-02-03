package peer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/jmoiron/sqlx"
)

type Peer struct {
	db     *sqlx.DB
	root   *resolvers.RootResolver
	schema *graphql.Schema
}

var _ api.Service = (*Peer)(nil)

type PeerConfig struct {
	VarDir      string
	GUIEndpoint string
}

func NewPeer(ctx context.Context, cfg PeerConfig) (*Peer, error) {
	// XXX reconsile SQL DB opening with exod/main.go
	dbPath := filepath.Join(cfg.VarDir, "exo.sqlite3")
	txMode := "exclusive"
	connStr := dbPath + "?_txlock=" + txMode
	db, err := sqlx.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db: %w", err)
	}
	r := &resolvers.RootResolver{
		DB:            db,
		SystemLog:     logging.CurrentLogger(ctx),
		ULIDGenerator: gensym.NewULIDGenerator(ctx),
		Routes:        resolvers.NewRoutesResolver(cfg.GUIEndpoint),
	}
	// XXX migration probably shouldn't happen here.
	if err := r.Migrate(ctx); err != nil {
		cmdutil.Fatalf("migrating db: %w", err)
	}
	return &Peer{
		db:     db,
		root:   r,
		schema: resolvers.NewSchema(r),
	}, nil
}

func (p *Peer) Do(ctx context.Context, doc string, vars map[string]interface{}, res interface{}) error {
	ctxVars := api.CurrentContextVariables(ctx)
	if ctxVars != nil && ctxVars.TaskID != "" {
		ctx = logging.ContextWithLogger(ctx, &api.EventLogger{
			Service:    p,
			SystemLog:  p.root.SystemLog,
			SourceType: "Task",
			SourceID:   ctxVars.TaskID,
		})
	}

	vars, _ = jsonutil.MustSimplify(vars).(map[string]interface{})
	resp := p.schema.Exec(ctx, doc, "", vars)
	for i, err := range resp.Errors {
		externalErr := sanitizeQueryError(err)
		if err == externalErr {
			continue
		}
		logging.Infof(ctx, "internal graphql error: %v", err)
		resp.Errors[i] = externalErr
	}
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

func sanitizeQueryError(original *gqlerrors.QueryError) *gqlerrors.QueryError {
	sanitized := gqlerrors.QueryError{
		Locations: original.Locations,
		Path:      original.Path,
	}
	if original.ResolverError != nil {
		sanitized.ResolverError = sanitizeError(original.ResolverError)
		if sanitized.ResolverError == original.ResolverError {
			return original
		}
		sanitized.Err = sanitized.ResolverError
	} else if original.Err != nil {
		sanitized.Err = sanitizeError(original.Err)
		if sanitized.Err == original.Err {
			return original
		}
	} else {
		sanitized.Err = sanitizeError(errors.New(original.Message))
	}
	sanitized.Message = sanitized.Err.Error()
	sanitized.Extensions = errorExtensions(sanitized.Err)
	return &sanitized
}

func sanitizeError(err error) error {
	if isExternalError(err) {
		return err
	}
	return errutil.InternalServerError
}

func isExternalError(err error) bool {
	_, ok := err.(errutil.HTTPError)
	return ok
}

func errorExtensions(err error) map[string]interface{} {
	iface, ok := err.(interface{ Extensions() map[string]interface{} })
	if !ok {
		return nil
	}
	return iface.Extensions()
}
