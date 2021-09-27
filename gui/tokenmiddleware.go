package gui

import (
	"net/http"
	"strings"
	"time"
)

type tokenCookieMiddleware struct {
	URL  string
	Next http.Handler
}

func (m *tokenCookieMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	if token == "" {
		m.Next.ServeHTTP(w, req)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		Secure:   !strings.HasPrefix(m.URL, "http:"),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})
	redirectURL := req.URL
	redirectQuery := redirectURL.Query()
	redirectQuery.Del("token")
	redirectURL.RawQuery = redirectQuery.Encode()
	http.Redirect(w, req, redirectURL.String(), http.StatusTemporaryRedirect)
}
