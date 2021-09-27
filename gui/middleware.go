package gui

import (
	"net/http"
	"strings"
	"time"
)

type guiMiddleware struct {
	URL  string
	Next http.Handler
}

func (m *guiMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Prevent click-jacking attacks.
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	// Strip token query parameter and set cookie.
	token := req.URL.Query().Get("token")
	if token != "" {
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
		return
	}

	m.Next.ServeHTTP(w, req)
}
