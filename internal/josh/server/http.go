package server

import (
	"net/http"
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
