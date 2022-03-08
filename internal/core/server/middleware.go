package server

import (
	"net/http"
)

type httpMiddleware interface {
	ServeHTTPMiddleware(w http.ResponseWriter, req *http.Request, next http.Handler)
}

type middlewareHandler struct {
	Middleware httpMiddleware
	Handler    http.Handler
}

func applyMiddleware(handler http.Handler, middleware ...httpMiddleware) http.Handler {
	h := handler
	for _, m := range middleware {
		h = &middlewareHandler{
			Middleware: m,
			Handler:    h,
		}
	}
	return h
}

func (m *middlewareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.Middleware.ServeHTTPMiddleware(w, req, m.Handler)
}
