// +build !bundle

package gui

import (
	"context"
	"fmt"
	"net/http"
	goutil "net/http/httputil"
	"net/url"
	"time"

	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/util/errutil"
	exoutil "github.com/deref/exo/internal/util/httputil"
)

var tokenCookieMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		token := query.Get("token")
		if token == "" {
			next.ServeHTTP(w, req)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().AddDate(1, 0, 0), // 1 year
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
		})
		query.Del("token")
		redirectURL := req.URL
		redirectURL.RawQuery = query.Encode()
		http.Redirect(w, req, redirectURL.String(), http.StatusTemporaryRedirect)
	})
}

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
