// +build !bundle

package gui

import (
	"context"
	"net/http"
	goutil "net/http/httputil"
	"net/url"

	"github.com/deref/exo/util/errutil"
	exoutil "github.com/deref/exo/util/httputil"
)

func NewHandler(ctx context.Context) http.Handler {
	guiURL, err := url.Parse("http://localhost:3000/")
	if err != nil {
		panic(err)
	}
	proxy := goutil.NewSingleHostReverseProxy(guiURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		err = errutil.NewHTTPError(http.StatusBadGateway, err.Error())
		exoutil.WriteError(w, req, err)
	}
	return proxy
}
