package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/deref/exo/util/errutil"
)

func handleOptions(w http.ResponseWriter, req *http.Request) bool {
	// TODO: Figure out the "right" thing to do for cors.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func writeError(w http.ResponseWriter, req *http.Request, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"
	if err, ok := err.(errutil.HTTPError); ok {
		status = err.HTTPStatus()
		message = err.Error()
	}
	if status == http.StatusInternalServerError {
		log.Printf("error processing request: %v", err)
	}
	writeJSON(w, status, map[string]interface{}{
		"status":  status,
		"message": message,
	})
}
