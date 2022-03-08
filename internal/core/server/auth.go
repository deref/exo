package server

import (
	"net/http"
	"strings"

	"github.com/deref/exo/internal/token"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/httputil"
)

type authMiddleware struct {
	TokenClient token.TokenClient
}

func (m *authMiddleware) ServeHTTPMiddleware(w http.ResponseWriter, req *http.Request, next http.Handler) {
	token := ""
	bearerSuffix := "Bearer "
	authHeader := req.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, bearerSuffix) {
		token = strings.TrimPrefix(authHeader, bearerSuffix)
	} else if cookie, err := req.Cookie("token"); err == nil {
		token = cookie.Value
	}

	authed, err := m.TokenClient.CheckToken(token)
	if err != nil {
		httputil.WriteError(w, req, errutil.NewHTTPError(http.StatusInternalServerError, "Could not validate token"))
		return
	}
	if !authed {
		httputil.WriteError(w, req, errutil.NewHTTPError(http.StatusUnauthorized, "Bad or no token"))
		return
	}
	next.ServeHTTP(w, req)
}
