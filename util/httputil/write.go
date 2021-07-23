package httputil

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func WriteString(w http.ResponseWriter, status int, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	io.WriteString(w, s)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
