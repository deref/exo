package josh

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

type MethodHandler struct {
	f   reflect.Value
	in  reflect.Type
	out reflect.Type
}

func NewMethodHandler(f interface{}) *MethodHandler {
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		panic("expected function, got " + fv.Kind().String())
	}
	typ := fv.Type()
	if typ.NumIn() != 2 {
		panic("wrong number of parameters. " + expectedSignature)
	}
	if typ.NumOut() != 2 {
		panic("wrong number of results. " + expectedSignature)
	}
	if typ.In(0) != contextType {
		panic("first parameter must be a context. " + expectedSignature)
	}
	in := typ.In(1)
	if in.Kind() != reflect.Ptr {
		panic("first parameter must be a pointer. " + expectedSignature)
	}
	out := typ.Out(0)
	if out.Kind() != reflect.Ptr {
		panic("first result must be a pointer. " + expectedSignature)
	}
	if typ.Out(1) != errorType {
		panic("second result must be an error. " + expectedSignature)
	}
	return &MethodHandler{
		f:   fv,
		in:  in.Elem(),
		out: out,
	}
}

const expectedSignature = "expected signature: func (ctx context.Context, input *YourInput) (*YourOutput, error)"

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (handler *MethodHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: Check content-type, accepts, etc.
	// TODO: Include a Request ID in logs?
	// TODO: Differentiate 400s from 500s, internal vs external errors, etc.
	dec := json.NewDecoder(req.Body)
	input := reflect.New(handler.in)
	if err := dec.Decode(input.Interface()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error parsing request: %v", err)
		w.Write([]byte("error parsing json\n"))
		return
	}
	results := handler.f.Call([]reflect.Value{
		reflect.ValueOf(req.Context()),
		input,
	})
	output := results[0].Interface()
	errv := results[1].Interface()
	if errv != nil {
		err := errv.(error)
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error processing request: %v", err)
		w.Write([]byte("internal server error\n"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(output); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
