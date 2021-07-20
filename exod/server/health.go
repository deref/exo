package server

import "net/http"

var HandleHealth http.HandlerFunc = handleHealth

// Note [HEALTH_CHECK]: This is the simplest possible health check so that
// the CLI can wait for the server to be listening. When this evolves, the CLI
// will also need to change.
func handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok\n"))
}
