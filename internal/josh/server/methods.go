package server

import (
	"context"
	"net/http"
	"path"
	"reflect"
	"time"

	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/httputil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
)

type MethodHandler struct {
	Factory func(req *http.Request) interface{}
	Name    string
}

const expectedSignature = "expected signature: func (ctx context.Context, input *YourInput) (*YourOutput, error)"

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (handler *MethodHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		err := errutil.HTTPErrorf(http.StatusMethodNotAllowed, "method not allowed: %q", req.Method)
		httputil.WriteError(w, req, err)
		return
	}

	ctx := req.Context()
	logger := logging.CurrentLogger(ctx)
	tel := telemetry.FromContext(ctx)

	start := time.Now()
	telOp := telemetry.OperationInvocation{
		Operation: handler.Name,
	}
	defer func() {
		telOp.DurationMicros = int(time.Since(start).Microseconds())
		tel.RecordOperation(telOp)
	}()

	method := reflect.ValueOf(handler.Factory(req))
	inputType := method.Type().In(1).Elem()

	// TODO: Check content-type, accepts, etc.
	// TODO: Include a Request ID in logs?
	// TODO: Differentiate 400s from 500s, internal vs external errors, etc.
	input := reflect.New(inputType)
	if req.Body != nil {
		if err := jsonutil.UnmarshalReader(req.Body, input.Interface()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Infof("error parsing request: %v", err)
			w.Write([]byte("error parsing json\n"))
			return
		}
	}
	results := method.Call([]reflect.Value{
		reflect.ValueOf(req.Context()),
		input,
	})
	output := results[0].Interface()
	errv := results[1].Interface()
	if errv != nil {
		err := errv.(error)
		httputil.WriteError(w, req, err)
		return
	}

	telOp.Success = true
	httputil.WriteJSON(w, req, http.StatusOK, output)
}

type MiddlewareMethod func(h http.Handler) http.Handler

type MuxBuilder struct {
	prefix     string
	mux        *http.ServeMux
	methods    []string
	middleware []MiddlewareMethod
}

func NewMuxBuilder(prefix string) *MuxBuilder {
	return &MuxBuilder{
		prefix: prefix,
		mux:    http.NewServeMux(),
	}
}

func (b *MuxBuilder) AddMiddleware(m MiddlewareMethod) {
	b.middleware = append(b.middleware, m)
}

func (b *MuxBuilder) Build() *http.ServeMux {
	b.end()
	mux := b.mux
	b.mux = nil
	return mux
}

func (b *MuxBuilder) AddMethod(name string, factory func(req *http.Request) interface{}) {
	joinedMiddleware := func(h http.Handler) http.Handler { return h } // Identity function
	for _, m := range b.middleware {
		joinedMiddleware = func(h http.Handler) http.Handler {
			return m(h)
		}
	}
	b.mux.Handle(path.Join(b.prefix, name), joinedMiddleware(&MethodHandler{
		Factory: factory,
		Name:    name,
	}))
	b.methods = append(b.methods, name)
}

func (b *MuxBuilder) Begin(prefix string) func() {
	oldPrefix := b.prefix
	oldMethods := b.methods
	b.prefix += prefix
	b.methods = nil
	return func() {
		b.end()
		b.prefix = oldPrefix
		b.methods = oldMethods
	}
}

func (b *MuxBuilder) end() {
	b.mux.Handle(b.prefix, &IntrospectionHandler{
		MethodNames: b.methods,
	})
}
