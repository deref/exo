package httputil

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/deref/exo/internal/util/logging"
)

func WriteString(w http.ResponseWriter, status int, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	io.WriteString(w, s)
}

func WriteJSON(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		logger := logging.CurrentLogger(req.Context())
		logger.Infof("error encoding response: %v", err)
	}
}
