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

func (p *Peer) Do(ctx context.Context, out interface{}, doc string, vars map[string]interface{}) error {
	ctx, operationName, vars := p.prepareOperation(ctx, vars)
	fmt.Printf("vars: %#v\n", vars)
	resp := p.schema.Exec(ctx, doc, operationName, vars)
	return p.handleResponse(ctx, out, resp)
}

func (p *Peer) Subscribe(ctx context.Context, out interface{}, doc string, vars map[string]interface{}) api.Subscription {
	ctx, operationName, vars := p.prepareOperation(ctx, vars)
	ctx, cancel := context.WithCancel(ctx)

	respC, err := p.schema.Subscribe(ctx, doc, operationName, vars)
	if err != nil {
		// Only schema configuration errors are checked synchronously.
		panic(err)
	}

	moreC := make(chan bool)
	errC := make(chan error, 1)
	go func() {
		defer close(moreC)
		defer cancel()
		for resp := range respC {
			if err := p.handleResponse(ctx, out, resp.(*graphql.Response)); err != nil {
				errC <- err
				return
			}
			select {
			case moreC <- true:
			case <-ctx.Done():
				return
			}
		}
	}()
	return &Subscription{
		moreC:  moreC,
		errC:   errC,
		cancel: cancel,
	}
}

func (p *Peer) prepareOperation(ctx context.Context, vars map[string]interface{}) (_ context.Context, operationName string, _ map[string]interface{}) {
	ctxVars := api.CurrentContextVariables(ctx)
	if ctxVars != nil && ctxVars.TaskID != "" {
		ctx = logging.ContextWithLogger(ctx, &api.EventLogger{
			Service:    p,
			SystemLog:  p.root.SystemLog,
			SourceType: "Task",
			SourceID:   ctxVars.TaskID,
		})
	}

	operationName = "" // TODO: Allow caller to specify somehow?

	if vars == nil {
		vars = make(map[string]interface{})
	} else {
		vars = jsonutil.MustSimplify(vars).(map[string]interface{})
	}

	return ctx, operationName, vars
}

func (p *Peer) handleResponse(ctx context.Context, out interface{}, resp *graphql.Response) error {
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
	if err := json.Unmarshal(resp.Data, out); err != nil {
		return err
	}
	return nil
}

type Subscription struct {
	moreC  <-chan bool
	errC   <-chan error
	err    error
	cancel func()
}

func (sub *Subscription) C() <-chan bool {
	return sub.moreC
}

func (sub *Subscription) Err() error {
	if sub.err == nil {
		select {
		case sub.err = <-sub.errC:
		default:
		}
	}
	return sub.err
}

func (sub *Subscription) Stop() {
	sub.cancel()
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
