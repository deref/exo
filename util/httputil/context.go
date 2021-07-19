package httputil

import (
	"context"
	"net/http"
)

func HandlerWithContext(ctx context.Context, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
