package peer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
)

type Peer struct {
	SystemLog   logging.Logger
	VarDir      string
	GUIEndpoint string
	Debug       bool

	root   *resolvers.RootResolver
	schema *graphql.Schema
}

var _ api.Service = (*Peer)(nil)

func (p *Peer) Init(ctx context.Context) error {
	p.root = &resolvers.RootResolver{
		VarDir:      p.VarDir,
		SystemLog:   p.SystemLog,
		GUIEndpoint: p.GUIEndpoint,
		Service:     p,
	}
	if err := p.root.Init(ctx); err != nil {
		return err
	}
	p.schema = resolvers.NewSchema(p.root)
	return nil
}

func (p *Peer) Do(ctx context.Context, out any, doc string, vars map[string]any) error {
	ctx, operationName, vars := p.prepareOperation(ctx, vars)
	resp := p.schema.Exec(ctx, doc, operationName, vars)
	return p.handleResponse(ctx, out, resp)
}

func (p *Peer) Subscribe(ctx context.Context, newRes func() any, doc string, vars map[string]any) api.Subscription {
	ctx, operationName, vars := p.prepareOperation(ctx, vars)
	ctx, cancel := context.WithCancel(ctx)

	source, err := p.schema.Subscribe(ctx, doc, operationName, vars)
	if err != nil {
		// Only schema configuration errors are checked synchronously.
		panic(err)
	}

	sink := make(chan any)
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

func (p *Peer) prepareOperation(ctx context.Context, vars map[string]any) (_ context.Context, operationName string, _ map[string]any) {
	ctxVars := api.CurrentContextVariables(ctx)
	if ctxVars != nil && ctxVars.TaskID != "" {
		ctx = logging.ContextWithLogger(ctx, &api.EventLogger{
			Service:    p,
			SystemLog:  p.SystemLog,
			SourceType: "Task",
			SourceID:   ctxVars.TaskID,
		})
	}

	operationName = "" // TODO: Allow caller to specify somehow?

	if vars == nil {
		vars = make(map[string]any)
	} else {
		vars = jsonutil.MustSimplify(vars).(map[string]any)
	}

	return ctx, operationName, vars
}

func (p *Peer) handleResponse(ctx context.Context, out any, resp *graphql.Response) error {
	for i, err := range resp.Errors {
		if errors.Is(err, context.Canceled) {
			// Unfortunately, there is no way to differentiate overall execution
			// errors from resolver errors. Least bad option is to treat cancelation
			// as coming from the overall execution, even though it may incorrectly
			// capture cancelations leaking out internal use within an individual
			// resolver.
			return err
		}
		externalErr := p.sanitizeQueryError(err)
		if err == externalErr {
			continue
		}
		logging.Infof(ctx, "internal graphql error: %v", err)
		resp.Errors[i] = externalErr
	}
	// NOTE [GRAPHQL_PARTIAL_FAILURE]: Returning on error before unmarshalling
	// promotes partial failure (GraphQL errors at specific paths) into total
	// failure. Not ideal for robustness, but greatly simplifies error handling
	// logic and acceptable for our not-so-distributed-systems use case.
	if len(resp.Errors) > 0 {
		return api.QueryErrorSet(resp.Errors)
	}
	return json.Unmarshal(resp.Data, out)
}

type Subscription struct {
	events <-chan any
	errs   <-chan error
	err    error
	cancel func()
}

func (sub *Subscription) Events() <-chan any {
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
	return p.root.Shutdown(ctx)
}

func (p *Peer) sanitizeQueryError(original *gqlerrors.QueryError) *gqlerrors.QueryError {
	if len(original.Locations) > 0 {
		return original
	}
	sanitized := gqlerrors.QueryError{
		Path: original.Path,
	}
	if original.ResolverError != nil {
		sanitized.ResolverError = p.sanitizeError(original.ResolverError)
		if sanitized.ResolverError == original.ResolverError {
			return original
		}
		sanitized.Err = sanitized.ResolverError
	} else if original.Err != nil {
		sanitized.Err = p.sanitizeError(original.Err)
		if sanitized.Err == original.Err {
			return original
		}
	} else {
		sanitized.Err = p.sanitizeError(errors.New(original.Message))
	}
	sanitized.Message = sanitized.Err.Error()
	sanitized.Extensions = errorExtensions(sanitized.Err)
	return &sanitized
}

func (p *Peer) sanitizeError(err error) error {
	if p.Debug || isExternalError(err) {
		return err
	}
	return errutil.InternalServerError
}

func isExternalError(err error) bool {
	_, ok := err.(errutil.HTTPError)
	return ok
}

func errorExtensions(err error) map[string]any {
	iface, ok := err.(interface{ Extensions() map[string]any })
	if !ok {
		return nil
	}
	return iface.Extensions()
}
