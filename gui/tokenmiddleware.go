package gui

import (
	"net/http"
	"time"
)

var tokenCookieMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := req.URL.Query().Get("token")
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
		redirectURL := req.URL
		redirectQuery := redirectURL.Query()
		redirectQuery.Del("token")
		redirectURL.RawQuery = redirectQuery.Encode()
		http.Redirect(w, req, redirectURL.String(), http.StatusTemporaryRedirect)
	})
}
