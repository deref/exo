package httputil

import (
	"net/http"

	"github.com/deref/exo/internal/util/errutil"
)

type HostAllowListHandler struct {
	Hosts []string
	Next  http.Handler
}

func (h *HostAllowListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, validHost := range h.Hosts {
		if req.Host == validHost {
			h.Next.ServeHTTP(w, req)
			return
		}
	}
	WriteError(w, req, errutil.NewHTTPError(http.StatusUnauthorized, "Invalid host header"))
}
