package server

import (
	"context"
	"net/http"
	"path"
	"reflect"

	"github.com/deref/exo/util/errutil"
	"github.com/deref/exo/util/httputil"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/logging"
)

type MethodHandler struct {
	Factory func(req *http.Request) interface{}
	Name    string
}

const expectedSignature = "expected signature: func (ctx context.Context, input *YourInput) (*YourOutput, error)"

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (handler *MethodHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handleOptions(w, req) {
		return
	}
	if req.Method != "POST" {
		err := errutil.HTTPErrorf(http.StatusMethodNotAllowed, "method not allowed: %q", req.Method)
		httputil.WriteError(w, req, err)
		return
	}

	logger := logging.CurrentLogger(req.Context())

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
	httputil.WriteJSON(w, req, http.StatusOK, output)
}

type MuxBuilder struct {
	prefix  string
	mux     *http.ServeMux
	methods []string
}

func NewMuxBuilder(prefix string) *MuxBuilder {
	return &MuxBuilder{
		prefix: prefix,
		mux:    http.NewServeMux(),
	}
}

func (b *MuxBuilder) Build() *http.ServeMux {
	b.end()
	mux := b.mux
	b.mux = nil
	return mux
}

func (b *MuxBuilder) AddMethod(name string, factory func(req *http.Request) interface{}) {
	b.mux.Handle(path.Join(b.prefix, name), &MethodHandler{
		Factory: factory,
		Name:    name,
	})
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
