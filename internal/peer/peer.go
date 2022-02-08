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
	resp := p.schema.Exec(ctx, doc, operationName, vars)
	return p.handleResponse(ctx, out, resp)
}

func (p *Peer) Subscribe(ctx context.Context, newRes func() interface{}, doc string, vars map[string]interface{}) api.Subscription {
	ctx, operationName, vars := p.prepareOperation(ctx, vars)
	ctx, cancel := context.WithCancel(ctx)

	source, err := p.schema.Subscribe(ctx, doc, operationName, vars)
	if err != nil {
		// Only schema configuration errors are checked synchronously.
		panic(err)
	}

	sink := make(chan interface{})
	errs := make(chan error, 1)
	go func() {
		defer close(sink)
		defer cancel()
		for resp := range source {
			out := newRes()
			resp := resp.(*graphql.Response)
			err := p.handleResponse(ctx, out, resp)
			if errors.Is(err, context.Canceled) {
				return
			}
			if err != nil {
				errs <- err
				return
			}
			select {
			case sink <- out:
			case <-ctx.Done():
				return
			}
		}
	}()
	return &Subscription{
		events: sink,
		errs:   errs,
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
		if errors.Is(err, context.Canceled) {
			// Unfortunately, there is no way to differentiate overall execution
			// errors from resolver errors. Least bad option is to treat cancelation
			// as coming from the overall execution, even though it may incorrectly
			// capture cancelations leaking out internal use within an individual
			// resolver.
			return err
		}
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
	return json.Unmarshal(resp.Data, out)
}

type Subscription struct {
	events <-chan interface{}
	errs   <-chan error
	err    error
	cancel func()
}

func (sub *Subscription) Events() <-chan interface{} {
	return sub.events
}

func (sub *Subscription) Err() error {
	if sub.err == nil {
		select {
		case sub.err = <-sub.errs:
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
	if len(original.Locations) > 0 {
		return original
	}
	sanitized := gqlerrors.QueryError{
		Path: original.Path,
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
