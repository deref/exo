// +build !bundle

package gui

import (
	"context"
	"fmt"
	"net/http"
	goutil "net/http/httputil"
	"net/url"

	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/util/errutil"
	exoutil "github.com/deref/exo/internal/util/httputil"
)

func NewHandler(ctx context.Context, cfg config.GUIConfig) http.Handler {
	guiURL, err := url.Parse(fmt.Sprintf("http://localhost:%d/", cfg.Port))
	if err != nil {
		panic(err)
	}
	proxy := goutil.NewSingleHostReverseProxy(guiURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		err = errutil.NewHTTPError(http.StatusBadGateway, err.Error())
		exoutil.WriteError(w, req, err)
	}
	return tokenCookieMiddleware(proxy)
}
