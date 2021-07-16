// +build !bundle

package gui

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewHandler(ctx context.Context) http.Handler {
	guiURL, err := url.Parse("http://localhost:3000/")
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(guiURL)
}
